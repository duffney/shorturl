package data

import "database/sql"

// Define a Models struct which wraps the MovieModel. We'll add other models to this, such as a UserModel and PermissionModel, as our build progresses.
type Models struct {
	Urls UrlModel
}

// Define a NewModels() function which returns a Models struct containing the initialized MovieModel.
func NewModels(db *sql.DB) Models {
	return Models{
		Urls: UrlModel{DB: db},
	}
}
