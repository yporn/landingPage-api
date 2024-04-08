package activityLogsRepositories

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/modules/activityLogs"
)

type IActivityLogsRepository interface {
	FindActivityLog() ([]activityLogs.ActivityLog, error)
}

type activityLogsRepository struct {
	db	*sqlx.DB
}

func ActivityLogsRepository(db *sqlx.DB) IActivityLogsRepository {
	return &activityLogsRepository{
		db:           db,
	}
}

func (r *activityLogsRepository) FindActivityLog() ([]activityLogs.ActivityLog, error) {
    query := `
        SELECT 
		"a"."id",
		"a"."user_id",
		"u"."name",
		"a"."action",
		"a"."details",
        "a"."created_at"
		FROM "activity_logs" "a"
		LEFT JOIN "users" AS "u" ON "a"."user_id" = "u"."id"
    `
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("query failed: %v", err)
    }
    defer rows.Close()

    var activityLogData []activityLogs.ActivityLog
    for rows.Next() {
        var activityLog activityLogs.ActivityLog
        // Initialize the User field before scanning into it
        activityLog.User = &activityLogs.User{}
        if err := rows.Scan(&activityLog.Id, &activityLog.User.Id, &activityLog.User.Name, &activityLog.Action, &activityLog.Detail, &activityLog.CreatedAt); err != nil {
            return nil, fmt.Errorf("scan activity log failed: %v", err)
        }
        activityLogData = append(activityLogData, activityLog)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("row error: %v", err)
    }

    return activityLogData, nil
}
