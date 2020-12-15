package persistence

import (
	"fmt"
	"fp-dynamic-elements-manager-controller/internal/logging/structs"
	structs2 "fp-dynamic-elements-manager-controller/internal/modules/structs"
	structs3 "fp-dynamic-elements-manager-controller/internal/queue/structs"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	ElementTypeTable = "module_element_type"
)

type ElementTypeRepo struct {
	db  *sqlx.DB
	log *structs.AppLogger
}

func NewElementTypeRepo(appDb *sqlx.DB, logger *structs.AppLogger) *ElementTypeRepo {
	return &ElementTypeRepo{db: appDb, log: logger}
}

func (e *ElementTypeRepo) InsertElementTypes(types structs2.ElementTypesWrapper, moduleId int64) {
	now := time.Now()

	smt := fmt.Sprintf(`INSERT INTO %s (created_at, updated_at, element_type, module_id) VALUES (?,?,?,?)`, ElementTypeTable)
	tx, err := e.db.Begin()

	if err != nil {
		e.log.SystemLogger.Error(err, "Error starting transaction to batch insert element types")
		return
	}

	for _, value := range types.ElementTypes {
		if e.exists(value, moduleId) {
			continue
		}
		_, err = tx.Exec(smt, now, now, value, moduleId)

		if err != nil {
			e.log.SystemLogger.Error(err, "Error batch inserting element types, rolling back")
			tx.Rollback()
			return
		}
	}

	err = tx.Commit()

	if err != nil {
		e.log.SystemLogger.Error(err, "Error committing batch insert element types")
		return
	}

	return
}

func (e *ElementTypeRepo) GetAllForModule(moduleId int64) (receiver []structs3.ElementType, err error) {
	err = e.db.Select(&receiver, fmt.Sprintf("SELECT element_type FROM %s WHERE module_id = ? ORDER BY created_at DESC;", ElementTypeTable), moduleId)
	return
}

func (e *ElementTypeRepo) exists(elementType structs3.ElementType, moduleId int64) bool {
	var eType structs2.ModuleElementType
	return e.db.Get(&eType, fmt.Sprintf("SELECT * FROM %s WHERE element_type = ? AND module_id = ? LIMIT 1;", ElementTypeTable), elementType, moduleId) == nil
}
