package slavemonitoring

import "database/sql"

type slaveStatus struct {
	SlaveIOState              sql.NullString `db:"slave_io_state" json:"slave_io_state"`
	MasterHost                sql.NullString `db:"master_host" json:"master_host"`
	MasterUser                sql.NullString `db:"Master_User" json:"master_user"`
	MasterPort                sql.NullString `db:"master_port" json:"master_port"`
	ConnectRetry              sql.NullString `db:"Connect_Retry" json:"connect_retry"`
	MasterLogFile             sql.NullString `db:"Master_Log_File" json:"master_log_file"`
	ReadMasterLogPos          sql.NullString `db:"Read_Master_Log_Pos" json:"read_master_log_pos"`
	RelayLogFile              sql.NullString `db:"Relay_Log_File" json:"relay_log_file"`
	RelayLogPos               sql.NullString `db:"Relay_Log_Pos" json:"relay_log_pos"`
	RelayMasterLogFile        sql.NullString `db:"Relay_Master_Log_File" json:"relay_master_log_file"`
	SlaveIORunning            sql.NullString `db:"Slave_IO_Running" json:"slave_io_running"`
	SlaveSQLRunning           sql.NullString `db:"Slave_SQL_Running" json:"slave_sql_running"`
	ReplicateDoDB             sql.NullString `db:"Replicate_Do_DB" json:"replicate_do_db"`
	ReplicateIgnoreDB         sql.NullString `db:"Replicate_Ignore_DB" json:"replicate_ignore_db"`
	ReplicateDoTable          sql.NullString `db:"Replicate_Do_Table" json:"replicate_do_table"`
	ReplicateIgnoreTable      sql.NullString `db:"Replicate_Ignore_Table" json:"replicate_ignore_table"`
	ReplicateWildDoTable      sql.NullString `db:"Replicate_Wild_Do_Table" json:"replicate_wild_do_table"`
	ReplicateWildIgnoreTable  sql.NullString `db:"Replicate_Wild_Ignore_Table" json:"replicate_wild_ignore_table"`
	LastErrno                 sql.NullInt64  `db:"Last_Errno" json:"last_errno"`
	LastError                 sql.NullString `db:"Last_Error" json:"last_error"`
	SkipCounter               sql.NullString `db:"Skip_Counter" json:"skip_counter"`
	ExecMasterLogPos          sql.NullString `db:"Exec_Master_Log_Pos" json:"exec_master_log_pos"`
	RelayLogSpace             sql.NullString `db:"Relay_Log_Space" json:"relay_log_space"`
	UntilCondition            sql.NullString `db:"Until_Condition" json:"until_condition"`
	UntilLogFile              sql.NullString `db:"Until_Log_File" json:"until_log_file"`
	UntilLogPos               sql.NullString `db:"Until_Log_Pos" json:"until_log_pos"`
	MasterSSLAllowed          sql.NullString `db:"Master_SSL_Allowed" json:"master_ssl_allowed"`
	MasterSSLCAFile           sql.NullString `db:"Master_SSL_CA_File" json:"master_sslca_file"`
	MasterSSLCAPath           sql.NullString `db:"Master_SSL_CA_Path" json:"master_sslca_path"`
	MasterSSLCert             sql.NullString `db:"Master_SSL_Cert" json:"master_ssl_cert"`
	MasterSSLCipher           sql.NullString `db:"Master_SSL_Cipher" json:"master_ssl_cipher"`
	MasterSSLKey              sql.NullString `db:"Master_SSL_Key" json:"master_ssl_key"`
	SecondsBehindMaster       sql.NullInt64  `db:"Seconds_Behind_Master" json:"seconds_behind_master"`
	MasterSSLVerifyServerCert sql.NullString `db:"Master_SSL_Verify_Server_Cert" json:"master_ssl_verify_server_cert"`
	LastIOErrno               sql.NullString `db:"Last_IO_Errno" json:"last_io_errno"`
	LastIOError               sql.NullString `db:"Last_IO_Error" json:"last_io_error"`
	LastSQLErrno              sql.NullString `db:"Last_SQL_Errno" json:"last_sql_errno"`
	LastSQLError              sql.NullString `db:"Last_SQL_Error" json:"last_sql_error"`
	ReplicateIgnoreServerIDs  sql.NullString `db:"Replicate_Ignore_Server_Ids" json:"replicate_ignore_server_i_ds"`
	MasterServerID            sql.NullString `db:"Master_Server_Id" json:"master_server_id"`
	MasterSSLCrl              sql.NullString `db:"Master_SSL_Crl" json:"master_ssl_crl"`
	MasterSSLCrlpath          sql.NullString `db:"Master_SSL_Crlpath" json:"master_ssl_crlpath"`
	UsingGtid                 sql.NullString `db:"Using_Gtid" json:"using_gtid"`
	GtidIOPos                 sql.NullString `db:"Gtid_IO_Pos" json:"gtid_io_pos"`
}
