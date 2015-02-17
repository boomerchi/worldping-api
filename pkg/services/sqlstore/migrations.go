package sqlstore

import . "github.com/grafana/grafana/pkg/services/sqlstore/migrator"

// --- Migration Guide line ---
// 1. Never change a migration that is committed and pushed to master
// 2. Always add new migrations (to change or undo previous migrations)
// 3. Some migraitons are not yet written (rename column, table, drop table, index etc)

func addMigrations(mg *Migrator) {
	addMigrationLogMigrations(mg)
	addUserMigrations(mg)
	addStarMigrations(mg)
	addAccountMigrations(mg)
	addDashboardMigration(mg)
	addDataSourceMigration(mg)
	addApiKeyMigrations(mg)
	addLocationMigrations(mg)
	addMonitorTypeMigrations(mg)
	addMonitorTypeSettingMigrations(mg)
	addMonitorMigrations(mg)
	addMonitorLocationMigrations(mg)
	addSiteMigrations(mg)

}

func addMigrationLogMigrations(mg *Migrator) {
	mg.AddMigration("create migration_log table", new(AddTableMigration).
		Name("migration_log").WithColumns(
		&Column{Name: "id", Type: DB_BigInt, IsPrimaryKey: true, IsAutoIncrement: true},
		&Column{Name: "migration_id", Type: DB_NVarchar, Length: 255},
		&Column{Name: "sql", Type: DB_Text},
		&Column{Name: "success", Type: DB_Bool},
		&Column{Name: "error", Type: DB_Text},
		&Column{Name: "timestamp", Type: DB_DateTime},
	))
}

func addUserMigrations(mg *Migrator) {
	mg.AddMigration("create user table", new(AddTableMigration).
		Name("user").WithColumns(
		&Column{Name: "id", Type: DB_BigInt, IsPrimaryKey: true, IsAutoIncrement: true},
		&Column{Name: "version", Type: DB_Int, Nullable: false},
		&Column{Name: "login", Type: DB_NVarchar, Length: 255, Nullable: false},
		&Column{Name: "email", Type: DB_NVarchar, Length: 255, Nullable: false},
		&Column{Name: "name", Type: DB_NVarchar, Length: 255, Nullable: true},
		&Column{Name: "password", Type: DB_NVarchar, Length: 255, Nullable: true},
		&Column{Name: "salt", Type: DB_NVarchar, Length: 50, Nullable: true},
		&Column{Name: "rands", Type: DB_NVarchar, Length: 50, Nullable: true},
		&Column{Name: "company", Type: DB_NVarchar, Length: 255, Nullable: true},
		&Column{Name: "account_id", Type: DB_BigInt, Nullable: false},
		&Column{Name: "is_admin", Type: DB_Bool, Nullable: false},
		&Column{Name: "created", Type: DB_DateTime, Nullable: false},
		&Column{Name: "updated", Type: DB_DateTime, Nullable: false},
	))

	mg.AddMigration("Add email_verified flag", new(AddColumnMigration).
		Table("user").Column(&Column{Name: "email_verified", Type: DB_Bool, Nullable: true}))

	mg.AddMigration("Add user.theme column", new(AddColumnMigration).
		Table("user").Column(&Column{Name: "theme", Type: DB_Varchar, Nullable: true, Length: 20}))

	//-------  user table indexes ------------------
	mg.AddMigration("add unique index user.login", new(AddIndexMigration).
		Table("user").Columns("login").Unique())
	mg.AddMigration("add unique index user.email", new(AddIndexMigration).
		Table("user").Columns("email").Unique())
}

func addStarMigrations(mg *Migrator) {
	mg.AddMigration("create star table", new(AddTableMigration).
		Name("star").WithColumns(
		&Column{Name: "id", Type: DB_BigInt, IsPrimaryKey: true, IsAutoIncrement: true},
		&Column{Name: "user_id", Type: DB_BigInt, Nullable: false},
		&Column{Name: "dashboard_id", Type: DB_BigInt, Nullable: false},
	))

	mg.AddMigration("add unique index star.user_id_dashboard_id", new(AddIndexMigration).
		Table("star").Columns("user_id", "dashboard_id").Unique())
}

func addAccountMigrations(mg *Migrator) {
	mg.AddMigration("create account table", new(AddTableMigration).
		Name("account").WithColumns(
		&Column{Name: "id", Type: DB_BigInt, IsPrimaryKey: true, IsAutoIncrement: true},
		&Column{Name: "version", Type: DB_Int, Nullable: false},
		&Column{Name: "name", Type: DB_NVarchar, Length: 255, Nullable: false},
		&Column{Name: "created", Type: DB_DateTime, Nullable: false},
		&Column{Name: "updated", Type: DB_DateTime, Nullable: false},
	))

	mg.AddMigration("add unique index account.name", new(AddIndexMigration).
		Table("account").Columns("name").Unique())

	//-------  account_user table -------------------
	mg.AddMigration("create account_user table", new(AddTableMigration).
		Name("account_user").WithColumns(
		&Column{Name: "id", Type: DB_BigInt, IsPrimaryKey: true, IsAutoIncrement: true},
		&Column{Name: "account_id", Type: DB_BigInt},
		&Column{Name: "user_id", Type: DB_BigInt},
		&Column{Name: "role", Type: DB_NVarchar, Length: 20},
		&Column{Name: "created", Type: DB_DateTime},
		&Column{Name: "updated", Type: DB_DateTime},
	))

	mg.AddMigration("add unique index account_user_aid_uid", new(AddIndexMigration).
		Name("aid_uid").Table("account_user").Columns("account_id", "user_id").Unique())
}

func addDashboardMigration(mg *Migrator) {
	mg.AddMigration("create dashboard table", new(AddTableMigration).
		Name("dashboard").WithColumns(
		&Column{Name: "id", Type: DB_BigInt, IsPrimaryKey: true, IsAutoIncrement: true},
		&Column{Name: "version", Type: DB_Int, Nullable: false},
		&Column{Name: "slug", Type: DB_NVarchar, Length: 255, Nullable: false},
		&Column{Name: "title", Type: DB_NVarchar, Length: 255, Nullable: false},
		&Column{Name: "data", Type: DB_Text, Nullable: false},
		&Column{Name: "account_id", Type: DB_BigInt, Nullable: false},
		&Column{Name: "created", Type: DB_DateTime, Nullable: false},
		&Column{Name: "updated", Type: DB_DateTime, Nullable: false},
	))

	mg.AddMigration("create dashboard_tag table", new(AddTableMigration).
		Name("dashboard_tag").WithColumns(
		&Column{Name: "id", Type: DB_BigInt, IsPrimaryKey: true, IsAutoIncrement: true},
		&Column{Name: "dashboard_id", Type: DB_BigInt, Nullable: false},
		&Column{Name: "term", Type: DB_NVarchar, Length: 50, Nullable: false},
	))

	//-------  indexes ------------------
	mg.AddMigration("add index dashboard.account_id", new(AddIndexMigration).
		Table("dashboard").Columns("account_id"))

	mg.AddMigration("add unique index dashboard_account_id_slug", new(AddIndexMigration).
		Table("dashboard").Columns("account_id", "slug").Unique())

	mg.AddMigration("add unique index dashboard_tag.dasboard_id_term", new(AddIndexMigration).
		Table("dashboard_tag").Columns("dashboard_id", "term").Unique())
}

func addDataSourceMigration(mg *Migrator) {
	mg.AddMigration("create data_source table", new(AddTableMigration).
		Name("data_source").WithColumns(
		&Column{Name: "id", Type: DB_BigInt, IsPrimaryKey: true, IsAutoIncrement: true},
		&Column{Name: "account_id", Type: DB_BigInt, Nullable: false},
		&Column{Name: "version", Type: DB_Int, Nullable: false},
		&Column{Name: "type", Type: DB_NVarchar, Length: 255, Nullable: false},
		&Column{Name: "name", Type: DB_NVarchar, Length: 255, Nullable: false},
		&Column{Name: "access", Type: DB_NVarchar, Length: 255, Nullable: false},
		&Column{Name: "url", Type: DB_NVarchar, Length: 255, Nullable: false},
		&Column{Name: "password", Type: DB_NVarchar, Length: 255, Nullable: true},
		&Column{Name: "user", Type: DB_NVarchar, Length: 255, Nullable: true},
		&Column{Name: "database", Type: DB_NVarchar, Length: 255, Nullable: true},
		&Column{Name: "basic_auth", Type: DB_Bool, Nullable: false},
		&Column{Name: "basic_auth_user", Type: DB_NVarchar, Length: 255, Nullable: true},
		&Column{Name: "basic_auth_password", Type: DB_NVarchar, Length: 255, Nullable: true},
		&Column{Name: "is_default", Type: DB_Bool, Nullable: false},
		&Column{Name: "created", Type: DB_DateTime, Nullable: false},
		&Column{Name: "updated", Type: DB_DateTime, Nullable: false},
	))

	//-------  indexes ------------------
	mg.AddMigration("add index data_source.account_id", new(AddIndexMigration).
		Table("data_source").Columns("account_id"))

	mg.AddMigration("add unique index data_source.account_id_name", new(AddIndexMigration).
		Table("data_source").Columns("account_id", "name").Unique())
}

func addApiKeyMigrations(mg *Migrator) {
	mg.AddMigration("create api_key table", new(AddTableMigration).
		Name("api_key").WithColumns(
		&Column{Name: "id", Type: DB_BigInt, IsPrimaryKey: true, IsAutoIncrement: true},
		&Column{Name: "account_id", Type: DB_BigInt, Nullable: false},
		&Column{Name: "name", Type: DB_NVarchar, Length: 255, Nullable: false},
		&Column{Name: "key", Type: DB_Varchar, Length: 64, Nullable: false},
		&Column{Name: "role", Type: DB_NVarchar, Length: 255, Nullable: false},
		&Column{Name: "created", Type: DB_DateTime, Nullable: false},
		&Column{Name: "updated", Type: DB_DateTime, Nullable: false},
	))

	//-------  indexes ------------------
	mg.AddMigration("add index api_key.account_id", new(AddIndexMigration).
		Table("api_key").Columns("account_id"))

	mg.AddMigration("add index api_key.key", new(AddIndexMigration).
		Table("api_key").Columns("key").Unique())

	mg.AddMigration("add index api_key.account_id_name", new(AddIndexMigration).
		Table("api_key").Columns("account_id", "name").Unique())
}

func addLocationMigrations(mg *Migrator) {
	mg.AddMigration("create location table", new(AddTableMigration).
		Name("location").WithColumns(
		&Column{Name: "id", Type: DB_BigInt, IsPrimaryKey: true, IsAutoIncrement: true},
		&Column{Name: "account_id", Type: DB_BigInt, Nullable: false},
		&Column{Name: "slug", Type: DB_NVarchar, Length: 255, Nullable: false},
		&Column{Name: "name", Type: DB_NVarchar, Length: 255, Nullable: false},
		&Column{Name: "country", Type: DB_NVarchar, Length: 255, Nullable: false},
		&Column{Name: "region", Type: DB_NVarchar, Length: 255, Nullable: false},
		&Column{Name: "provider", Type: DB_NVarchar, Length: 255, Nullable: false},
		&Column{Name: "public", Type: DB_Bool, Nullable: false},
		&Column{Name: "created", Type: DB_DateTime, Nullable: false},
		&Column{Name: "updated", Type: DB_DateTime, Nullable: false},
	))

	//-------  indexes ------------------

	mg.AddMigration("add unique index location.account_id_slug", new(AddIndexMigration).
		Table("location").Columns("account_id", "slug").Unique())
}

func addMonitorTypeMigrations(mg *Migrator) {
	mg.AddMigration("create monitor_type table", new(AddTableMigration).
		Name("monitor_type").WithColumns(
		&Column{Name: "id", Type: DB_BigInt, IsPrimaryKey: true, IsAutoIncrement: true},
		&Column{Name: "name", Type: DB_NVarchar, Length: 255, Nullable: false},
		&Column{Name: "created", Type: DB_DateTime, Nullable: false},
		&Column{Name: "updated", Type: DB_DateTime, Nullable: false},
	))

	//-------  data ------------------
	mg.AddMigration("insert http type into monitor_type table", new(RawSqlMigration).
		Sqlite("INSERT INTO monitor_type values(1,'HTTP',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)").
		Mysql("INSERT INTO monitor_type values(1,'HTTP',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)"))
	mg.AddMigration("insert https type into monitor_type table", new(RawSqlMigration).
		Sqlite("INSERT INTO monitor_type values(2,'HTTPS',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)").
		Mysql("INSERT INTO monitor_type values(2,'HTTPS',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)"))
	mg.AddMigration("insert ping type into monitor_type table", new(RawSqlMigration).
		Sqlite("INSERT INTO monitor_type values(3,'Ping',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)").
		Mysql("INSERT INTO monitor_type values(3,'Ping',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)"))
	mg.AddMigration("insert dns type into monitor_type table", new(RawSqlMigration).
		Sqlite("INSERT INTO monitor_type values(4,'DNS',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)").
		Mysql("INSERT INTO monitor_type values(4,'DNS',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)"))

}

func addMonitorTypeSettingMigrations(mg *Migrator) {
	mg.AddMigration("create monitor_type_setting table", new(AddTableMigration).
		Name("monitor_type_setting").WithColumns(
		&Column{Name: "id", Type: DB_BigInt, IsPrimaryKey: true, IsAutoIncrement: true},
		&Column{Name: "monitor_type_id", Type: DB_BigInt, Nullable: false},
		&Column{Name: "variable", Type: DB_NVarchar, Length: 255, Nullable: false},
		&Column{Name: "description", Type: DB_NVarchar, Length: 255, Nullable: true},
		&Column{Name: "data_type", Type: DB_NVarchar, Length: 255, Nullable: false},
		&Column{Name: "conditions", Type: DB_NVarchar, Length: 1024, Nullable: false},
		&Column{Name: "default_value", Type: DB_NVarchar, Length: 255, Nullable: false},
		&Column{Name: "required", Type: DB_Bool, Nullable: false},
		&Column{Name: "created", Type: DB_DateTime, Nullable: false},
		&Column{Name: "updated", Type: DB_DateTime, Nullable: false},
	))

	//-------  indexes ------------------
	mg.AddMigration("add index monitor_type_setting.monitor_type_id", new(AddIndexMigration).
		Table("monitor_type_setting").Columns("monitor_type_id"))

	//-------  data ------------------
	mg.AddMigration("insert http.host type_settings into monitor_type_setting table", new(RawSqlMigration).
		Sqlite("INSERT INTO monitor_type_setting values(null,1,'host','Hostname','String','{}','',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)").
		Mysql("INSERT INTO monitor_type_setting values(null,1,'host','Hostname','String','{}','',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)"))
	mg.AddMigration("insert http.path type_settings into monitor_type_setting table", new(RawSqlMigration).
		Sqlite("INSERT INTO monitor_type_setting values(null,1,'path','Path','String','{}','/',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)").
		Mysql("INSERT INTO monitor_type_setting values(null,1,'path','Path','String','{}','/',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)"))
	mg.AddMigration("insert http.port type_settings into monitor_type_setting table", new(RawSqlMigration).
		Sqlite("INSERT INTO monitor_type_setting values(null,1,'port','Port','Number','{}','80',0,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)").
		Mysql("INSERT INTO monitor_type_setting values(null,1,'port','Port','Number','{}','80',0,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)"))
	mg.AddMigration("insert http.method type_settings into monitor_type_setting table", new(RawSqlMigration).
		Sqlite("INSERT INTO monitor_type_setting values(null,1,'method','Method','Enum','{\"values\": [\"GET\", \"POST\",\"PUT\",\"DELETE\", \"HEAD\"]}','GET',0,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)").
		Mysql("INSERT INTO monitor_type_setting values(null,1,'method','Method','Enum','{\"values\": [\"GET\", \"POST\",\"PUT\",\"DELETE\", \"HEAD\"]}','GET',0,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)"))
	mg.AddMigration("insert http.headers type_settings into monitor_type_setting table", new(RawSqlMigration).
		Sqlite("INSERT INTO monitor_type_setting values(null,1,'headers','Headers','Text','{}','Accept-Encoding: gzip\nUser-Agent: raintank collector\n',0,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)").
		Mysql("INSERT INTO monitor_type_setting values(null,1,'headers','Headers','Text','{}','Accept-Encoding: gzip\nUser-Agent: raintank collector\n',0,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)"))
	mg.AddMigration("insert http.expectRegex type_settings into monitor_type_setting table", new(RawSqlMigration).
		Sqlite("INSERT INTO monitor_type_setting values(null,1,'expectRegex','Content Match','String','{}','',0,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)").
		Mysql("INSERT INTO monitor_type_setting values(null,1,'expectRegex','Content Match','String','{}','',0,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)"))

	mg.AddMigration("insert https.host type_settings into monitor_type_setting table", new(RawSqlMigration).
		Sqlite("INSERT INTO monitor_type_setting values(null,2,'host','Hostname','String','{}','',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)").
		Mysql("INSERT INTO monitor_type_setting values(null,2,'host','Hostname','String','{}','',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)"))
	mg.AddMigration("insert https.path type_settings into monitor_type_setting table", new(RawSqlMigration).
		Sqlite("INSERT INTO monitor_type_setting values(null,2,'path','Path','String','{}','/',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)").
		Mysql("INSERT INTO monitor_type_setting values(null,2,'path','Path','String','{}','/',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)"))
	mg.AddMigration("insert https.validateCert type_settings into monitor_type_setting table", new(RawSqlMigration).
		Sqlite("INSERT INTO monitor_type_setting values(null,2,'validateCert','Validate SSL Certificate','Boolean','{}','true',0,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)").
		Mysql("INSERT INTO monitor_type_setting values(null,2,'validateCert','Validate SSL Certificate','Boolean','{}','true',0,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)"))
	mg.AddMigration("insert https.port type_settings into monitor_type_setting table", new(RawSqlMigration).
		Sqlite("INSERT INTO monitor_type_setting values(null,2,'port','Port','Number','{}','80',0,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)").
		Mysql("INSERT INTO monitor_type_setting values(null,2,'port','Port','Number','{}','80',0,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)"))
	mg.AddMigration("insert https.method type_settings into monitor_type_setting table", new(RawSqlMigration).
		Sqlite("INSERT INTO monitor_type_setting values(null,2,'method','Method','Enum','{\"values\": [\"GET\", \"POST\",\"PUT\",\"DELETE\", \"HEAD\"]}','GET',0,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)").
		Mysql("INSERT INTO monitor_type_setting values(null,2,'method','Method','Enum','{\"values\": [\"GET\", \"POST\",\"PUT\",\"DELETE\", \"HEAD\"]}','GET',0,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)"))
	mg.AddMigration("insert https.headers type_settings into monitor_type_setting table", new(RawSqlMigration).
		Sqlite("INSERT INTO monitor_type_setting values(null,2,'headers','Headers','Text','{}','Accept-Encoding: gzip\nUser-Agent: raintank collector\n',0,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)").
		Mysql("INSERT INTO monitor_type_setting values(null,2,'headers','Headers','Text','{}','Accept-Encoding: gzip\nUser-Agent: raintank collector\n',0,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)"))
	mg.AddMigration("insert https.expectRegex type_settings into monitor_type_setting table", new(RawSqlMigration).
		Sqlite("INSERT INTO monitor_type_setting values(null,2,'expectRegex','Content Match','String','{}','',0,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)").
		Mysql("INSERT INTO monitor_type_setting values(null,2,'expectRegex','Content Match','String','{}','',0,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)"))

	mg.AddMigration("insert ping.hostname type_settings into monitor_type_setting table", new(RawSqlMigration).
		Sqlite("INSERT INTO monitor_type_setting values(null,3,'hostname','Hostname','String','{}','',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)").
		Mysql("INSERT INTO monitor_type_setting values(null,3,'hostname','Hostname','String','{}','',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)"))

	mg.AddMigration("insert dns.name type_settings into monitor_type_setting table", new(RawSqlMigration).
		Sqlite("INSERT INTO monitor_type_setting values(null,4,'name','Record Name','String','{}','',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)").
		Mysql("INSERT INTO monitor_type_setting values(null,4,'name','Record Name','String','{}','',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)"))
	mg.AddMigration("insert dns.type type_settings into monitor_type_setting table", new(RawSqlMigration).
		Sqlite("INSERT INTO monitor_type_setting values(null,4,'type','Record Tyoe','Enum','{\"values\": [\"A\",\"AAAA\",\"CNAME\",\"MX\",\"NS\",\"PTR\",\"SOA\",\"SRV\",\"TXT\"]}','A',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)").
		Mysql("INSERT INTO monitor_type_setting values(null,4,'type','Record Type','Enum','{\"values\": [\"A\",\"AAAA\",\"CNAME\",\"MX\",\"NS\",\"PTR\",\"SOA\",\"SRV\",\"TXT\"]}','A',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)"))
	mg.AddMigration("insert dns.server type_settings into monitor_type_setting table", new(RawSqlMigration).
		Sqlite("INSERT INTO monitor_type_setting values(null,4,'server','Server','String','{}','8.8.8.8',0,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)").
		Mysql("INSERT INTO monitor_type_setting values(null,4,'server','Server','String','{}','8.8.8.8',0,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)"))
	mg.AddMigration("insert dns.port type_settings into monitor_type_setting table", new(RawSqlMigration).
		Sqlite("INSERT INTO monitor_type_setting values(null,4,'port','port','Number','{}','53',0,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)").
		Mysql("INSERT INTO monitor_type_setting values(null,4,'port','Port','Number','{}','53',0,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)"))
	mg.AddMigration("insert dns.protocol type_settings into monitor_type_setting table", new(RawSqlMigration).
		Sqlite("INSERT INTO monitor_type_setting values(null,4,'protocol','Protocol','Enum','{\"values\": [\"tcp\",\"udp\"]}','udp',0,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)").
		Mysql("INSERT INTO monitor_type_setting values(null,4,'protocol','Protocol','Enum','{\"values\": [\"tcp\",\"udp\"]}','udp',0,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)"))

}

func addMonitorMigrations(mg *Migrator) {
	mg.AddMigration("create monitor table", new(AddTableMigration).
		Name("monitor").WithColumns(
		&Column{Name: "id", Type: DB_BigInt, IsPrimaryKey: true, IsAutoIncrement: true},
		&Column{Name: "account_id", Type: DB_BigInt, Nullable: false},
		&Column{Name: "slug", Type: DB_NVarchar, Length: 255, Nullable: false},
		&Column{Name: "name", Type: DB_NVarchar, Length: 255, Nullable: false},
		&Column{Name: "monitor_type_id", Type: DB_BigInt, Nullable: false},
		&Column{Name: "offset", Type: DB_BigInt, Nullable: false},
		&Column{Name: "frequency", Type: DB_BigInt, Nullable: false},
		&Column{Name: "enabled", Type: DB_Bool, Nullable: false},
		&Column{Name: "settings", Type: DB_NVarchar, Length: 2048, Nullable: false},
		&Column{Name: "created", Type: DB_DateTime, Nullable: false},
		&Column{Name: "updated", Type: DB_DateTime, Nullable: false},
	))

	//-------  indexes ------------------
	mg.AddMigration("add index monitor.monitor_type_id", new(AddIndexMigration).
		Table("monitor").Columns("monitor_type_id"))

	mg.AddMigration("add unique index monitor.account_id_slug", new(AddIndexMigration).
		Table("monitor").Columns("account_id", "slug").Unique())

	//------ new columns ----------------
	mg.AddMigration("add siteId column to monitor table", new(AddColumnMigration).
		Table("monitor").Column(&Column{Name: "site_id", Type: DB_BigInt, Nullable: true}))

	mg.AddMigration("add namespace column to monitor table", new(AddColumnMigration).
		Table("monitor").Column(&Column{Name: "namespace", Type: DB_NVarchar, Length: 255, Nullable: true}))

}

func addMonitorLocationMigrations(mg *Migrator) {
	mg.AddMigration("create monitor_location table", new(AddTableMigration).
		Name("monitor_location").WithColumns(
		&Column{Name: "id", Type: DB_BigInt, IsPrimaryKey: true, IsAutoIncrement: true},
		&Column{Name: "monitor_id", Type: DB_BigInt, Nullable: false},
		&Column{Name: "location_id", Type: DB_BigInt, Nullable: false},
	))

	//-------  indexes ------------------
	mg.AddMigration("add index monitor_location.monitor_id_location_id", new(AddIndexMigration).
		Table("monitor_location").Columns("monitor_id", "location_id"))
}

func addSiteMigrations(mg *Migrator) {
	mg.AddMigration("create site table", new(AddTableMigration).
		Name("site").WithColumns(
		&Column{Name: "id", Type: DB_BigInt, IsPrimaryKey: true, IsAutoIncrement: true},
		&Column{Name: "account_id", Type: DB_BigInt, Nullable: false},
		&Column{Name: "slug", Type: DB_NVarchar, Length: 255, Nullable: false},
		&Column{Name: "name", Type: DB_NVarchar, Length: 255, Nullable: false},
		&Column{Name: "created", Type: DB_DateTime, Nullable: false},
		&Column{Name: "updated", Type: DB_DateTime, Nullable: false},
	))

	//-------  indexes ------------------
	mg.AddMigration("add unique index site.account_id_slug", new(AddIndexMigration).
		Table("site").Columns("account_id", "slug").Unique(),
	)
}
