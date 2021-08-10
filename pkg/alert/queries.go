package alert

var upsertCriticalSlackCredentialQuery = `INSERT INTO slack_credentials (created_at,updated_at,channel_name, 
	level,team_name,entity) VALUES (now(), now(), @channel_name, 'CRITICAL', @team_name,
	@entity) ON CONFLICT (level, team_name) DO UPDATE SET "updated_at"= now(),"deleted_at"="excluded"."deleted_at",
	"channel_name"="excluded"."channel_name",
	"level"="excluded"."level", "team_name"="excluded"."team_name","entity"="excluded"."entity" RETURNING "id"`

var upsertWarningSlackCredentialQuery = `INSERT INTO slack_credentials (created_at,updated_at,channel_name, 
	level,team_name,entity) VALUES (now(), now(), @channel_name, 'WARNING', @team_name,
	@entity) ON CONFLICT (level, team_name) DO UPDATE SET "updated_at"= now(),"deleted_at"="excluded"."deleted_at",
	"channel_name"="excluded"."channel_name", 
	"level"="excluded"."level", "team_name"="excluded"."team_name","entity"="excluded"."entity" RETURNING "id"`

var upsertPagerdutyCredentialsQuery = `INSERT INTO pagerduty_credentials (created_at, updated_at, service_key,
	team_name, entity) VALUES(now(), now(), @service_key, @team_name, @entity)
	ON CONFLICT(team_name) DO UPDATE SET "updated_at" = now(), service_key = excluded.service_key, 
	entity = excluded.entity`

var joinQuery = `select sw.team_name as team_name, pg.service_key as pg_service_key, sw.channel_name as 
	warning_channel_name, sc.channel_name as critical_channel_name from slack_credentials as 
	sw join slack_credentials as sc on sc.team_name=sw.team_name join pagerduty_credentials as pg on sc.team_name = 
	pg.team_name where sw.entity= ? and sw.level='WARNING' and sc.entity=? and sc.level ='CRITICAL' and pg.entity=?`

var selectQuery = `select sw.entity as entity, pg.service_key as pg_service_key, sw.channel_name  
	as warning_channel_name, sc.channel_name as critical_channel_name from slack_credentials as sw join 
	slack_credentials as sc on sc.team_name=sw.team_name join pagerduty_credentials as pg on sc.team_name = pg.team_name
 	where sw.team_name= ? and sw.level='WARNING' and sc.level ='CRITICAL'`
