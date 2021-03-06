package auth


import (
    "errors"
    "runtime"

    "github.com/casbin/casbin/model"
    "github.com/casbin/casbin/persist"
    "github.com/jinzhu/gorm"
    obj_model "github.com/LongMarch7/go-web/db/object/models"
)


// Adapter represents the Gorm adapter for policy storage.
type Adapter struct {
    driverName     string
    dataSourceName string
    dbSpecified    bool
    db             *gorm.DB
}

// finalizer is the destructor for Adapter.
func finalizer(a *Adapter) {
    a.db.Close()
}

// NewAdapter is the constructor for Adapter.
// dbSpecified is an optional bool parameter. The default value is false.
// It's up to whether you have specified an existing DB in dataSourceName.
// If dbSpecified == true, you need to make sure the DB in dataSourceName exists.
// If dbSpecified == false, the adapter will automatically create a DB named "casbin".
func NewAdapter(driverName string, dataSourceName string, tableName string, dbSpecified ...bool) *Adapter {
    a := &Adapter{}
    a.driverName = driverName
    a.dataSourceName = dataSourceName

    if len(dbSpecified) == 0 {
        a.dbSpecified = false
    } else if len(dbSpecified) == 1 {
        a.dbSpecified = dbSpecified[0]
    } else {
        panic(errors.New("invalid parameter: dbSpecified"))
    }

    // Open the DB, create it if not existed.
    a.open(tableName)

    // Call the destructor when the object is released.
    runtime.SetFinalizer(a, finalizer)

    return a
}

func NewAdapterByDB(db *gorm.DB) *Adapter {
    a := &Adapter{
        db: db,
    }
    a.createTable()
    return a
}

//func (a *Adapter) createDatabase() error {
//    var err error
//    var db *gorm.DB
//    if a.driverName == "postgres" {
//        db, err = gorm.Open(a.driverName, a.dataSourceName+" dbname=postgres")
//    } else {
//        db, err = gorm.Open(a.driverName, a.dataSourceName)
//    }
//    if err != nil {
//        return err
//    }
//    defer db.Close()
//
//    if a.driverName == "postgres" {
//        if err = db.Exec("CREATE DATABASE casbin").Error; err != nil {
//            // 42P04 is	duplicate_database
//            if err.(*pq.Error).Code == "42P04" {
//                return nil
//            }
//        }
//    } else if a.driverName != "sqlite3" {
//        err = db.Exec("CREATE DATABASE IF NOT EXISTS casbin").Error
//    }
//    return err
//}

func (a *Adapter) open(tableName string) {
    var err error
    var db *gorm.DB

    if a.dbSpecified {
        db, err = gorm.Open(a.driverName, a.dataSourceName)
        if err != nil {
            panic(err)
        }
    } else {
        //if err = a.createDatabase(); err != nil {
        //    panic(err)
        //}

        if a.driverName == "postgres" {
            db, err = gorm.Open(a.driverName, a.dataSourceName + " dbname=" + tableName)
        } else if a.driverName == "sqlite3" {
            db, err = gorm.Open(a.driverName, a.dataSourceName)
        } else {
            db, err = gorm.Open(a.driverName, a.dataSourceName+tableName)
        }
        if err != nil {
            panic(err)
        }
    }

    a.db = db

    a.createTable()
}

func (a *Adapter) close() {
    a.db.Close()
    a.db = nil
}

func (a *Adapter) createTable() {
    if a.db.HasTable(&obj_model.MicroCasbinRule{}) {
        return
    }

    err := a.db.CreateTable(&obj_model.MicroCasbinRule{}).Error
    if err != nil {
        panic(err)
    }
}

func (a *Adapter) dropTable() {
    err := a.db.DropTable(&obj_model.MicroCasbinRule{}).Error
    if err != nil {
        panic(err)
    }
}

func loadPolicyLine(line obj_model.MicroCasbinRule, model model.Model) {
    lineText := line.PType
    if line.V0 != "" {
        lineText += ", " + line.V0
    }
    if line.V1 != "" {
        lineText += ", " + line.V1
    }
    if line.V2 != "" {
        lineText += ", " + line.V2
    }
    if line.V3 != "" {
        lineText += ", " + line.V3
    }
    if line.V4 != "" {
        lineText += ", " + line.V4
    }
    if line.V5 != "" {
        lineText += ", " + line.V5
    }

    persist.LoadPolicyLine(lineText, model)
}

// LoadPolicy loads policy from database.
func (a *Adapter) LoadPolicy(model model.Model) error {
    var lines []obj_model.MicroCasbinRule
    err := a.db.Find(&lines).Error
    if err != nil {
        return err
    }

    for _, line := range lines {
        loadPolicyLine(line, model)
    }

    return nil
}

func savePolicyLine(ptype string, rule []string) obj_model.MicroCasbinRule {
    line := obj_model.MicroCasbinRule{}

    line.PType = ptype
    if len(rule) > 0 {
        line.V0 = rule[0]
    }
    if len(rule) > 1 {
        line.V1 = rule[1]
    }
    if len(rule) > 2 {
        line.V2 = rule[2]
    }
    if len(rule) > 3 {
        line.V3 = rule[3]
    }
    if len(rule) > 4 {
        line.V4 = rule[4]
    }
    if len(rule) > 5 {
        line.V5 = rule[5]
    }

    return line
}

// SavePolicy saves policy to database.
func (a *Adapter) SavePolicy(model model.Model) error {
    a.dropTable()
    a.createTable()

    for ptype, ast := range model["p"] {
        for _, rule := range ast.Policy {
            line := savePolicyLine(ptype, rule)
            err := a.db.Create(&line).Error
            if err != nil {
                return err
            }
        }
    }

    for ptype, ast := range model["g"] {
        for _, rule := range ast.Policy {
            line := savePolicyLine(ptype, rule)
            err := a.db.Create(&line).Error
            if err != nil {
                return err
            }
        }
    }

    return nil
}

// AddPolicy adds a policy rule to the storage.
func (a *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
    line := savePolicyLine(ptype, rule)
    err := a.db.Create(&line).Error
    return err
}

// RemovePolicy removes a policy rule from the storage.
func (a *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
    line := savePolicyLine(ptype, rule)
    err := rawDelete(a.db, line) //can't use db.Delete as we're not using primary key http://jinzhu.me/gorm/crud.html#delete
    return err
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
    line := obj_model.MicroCasbinRule{}

    line.PType = ptype
    if fieldIndex <= 0 && 0 < fieldIndex + len(fieldValues) {
        line.V0 = fieldValues[0 - fieldIndex]
    }
    if fieldIndex <= 1 && 1 < fieldIndex + len(fieldValues) {
        line.V1 = fieldValues[1 - fieldIndex]
    }
    if fieldIndex <= 2 && 2 < fieldIndex + len(fieldValues) {
        line.V2 = fieldValues[2 - fieldIndex]
    }
    if fieldIndex <= 3 && 3 < fieldIndex + len(fieldValues) {
        line.V3 = fieldValues[3 - fieldIndex]
    }
    if fieldIndex <= 4 && 4 < fieldIndex + len(fieldValues) {
        line.V4 = fieldValues[4 - fieldIndex]
    }
    if fieldIndex <= 5 && 5 < fieldIndex + len(fieldValues) {
        line.V5 = fieldValues[5 - fieldIndex]
    }
    err := rawDelete(a.db, line)
    return err
}

func rawDelete(db *gorm.DB, line obj_model.MicroCasbinRule) error {
    queryArgs := []interface{}{line.PType}

    queryStr := "p_type = ?"
    if line.V0 != "" {
        queryStr += " and v0 = ?"
        queryArgs = append(queryArgs, line.V0)
    }
    if line.V1 != "" {
        queryStr += " and v1 = ?"
        queryArgs = append(queryArgs, line.V1)
    }
    if line.V2 != "" {
        queryStr += " and v2 = ?"
        queryArgs = append(queryArgs, line.V2)
    }
    if line.V3 != "" {
        queryStr += " and v3 = ?"
        queryArgs = append(queryArgs, line.V3)
    }
    if line.V4 != "" {
        queryStr += " and v4 = ?"
        queryArgs = append(queryArgs, line.V4)
    }
    if line.V5 != "" {
        queryStr += " and v5 = ?"
        queryArgs = append(queryArgs, line.V5)
    }
    args := append([]interface{}{queryStr}, queryArgs...)
    err := db.Delete(obj_model.MicroCasbinRule{}, args...).Error
    return err
}