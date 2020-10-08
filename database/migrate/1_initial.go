package migrate

import (
	"fmt"

	"github.com/go-pg/migrations/v8"
)

const uuidExtension = `
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
`

const adminTable = `
CREATE TABLE IF NOT EXISTS admins (
	id uuid DEFAULT uuid_generate_v4(),
	is_super_admin BOOLEAN NOT NULL DEFAULT FALSE,
	name VARCHAR(255) NOT NULL,
	email VARCHAR(255) NOT NULL,
	password varchar(255) NOT NULL,
	create_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	update_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY(id)
);
`

const projectTable = `
CREATE TABLE IF NOT EXISTS projects (
	id uuid DEFAULT uuid_generate_v4(),
	name text NOT NULL,
	admin_id uuid NOT NULL,
	file_name_pattern text NOT NULL,
	file_name_extensions text[] NOT NULL,
	file_name_example text NOT NULL,
	usable boolean NOT NULL DEFAULT TRUE,
	create_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	update_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id),
	FOREIGN KEY (admin_id) REFERENCES admins(id)
);
`

const submissionTable = `
CREATE TABLE IF NOT EXISTS submissions (
	id uuid DEFAULT uuid_generate_v4(),
	project_id uuid NOT NULL,
	secret text NOT NULL,
	file_name text NOT NULL,
	file_path text NOT NULL,
	md5 text NOT NULL,
	create_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	update_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	ip text,
	PRIMARY KEY (id),
	FOREIGN KEY (project_id) REFERENCES projects(id)
);
`

const invitationCodeTable = `
CREATE TABLE IF NOT EXISTS invitation_codes (
	id uuid DEFAULT uuid_generate_v4(),
	email text NOT NULL,
	code text NOT NULL,
	create_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	update_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY(id)
);
`

func init() {
	up := []string{
		uuidExtension,
		adminTable,
		projectTable,
		submissionTable,
		invitationCodeTable,
	}

	down := []string{
		"DROP TABLE IF EXISTS submissions",
		"DROP TABLE IF EXISTS projects",
		"DROP TABLE IF EXISTS admins",
		"DROP TABLE IF EXISTS invitation_codes",
	}

	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("creating initial tables")
		for _,query := range up {
		_,err := db.Exec(query);
			if err !=nil {
				return err
			}
		}
		return nil
	},func(db migrations.DB) error {
		fmt.Println("dropping initial tables")
		for _,query := range down {
			_,err := db.Exec(query)
			if err != nil {
				return err
			}
		}
		return nil
	})

}
