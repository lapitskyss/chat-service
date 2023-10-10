// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// ChatsColumns holds the columns for the "chats" table.
	ChatsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "client_id", Type: field.TypeUUID, Unique: true},
		{Name: "created_at", Type: field.TypeTime},
	}
	// ChatsTable holds the schema information for the "chats" table.
	ChatsTable = &schema.Table{
		Name:       "chats",
		Columns:    ChatsColumns,
		PrimaryKey: []*schema.Column{ChatsColumns[0]},
	}
	// MessagesColumns holds the columns for the "messages" table.
	MessagesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "author_id", Type: field.TypeUUID},
		{Name: "is_visible_for_client", Type: field.TypeBool},
		{Name: "is_visible_for_manager", Type: field.TypeBool},
		{Name: "body", Type: field.TypeString, Size: 4000},
		{Name: "checked_at", Type: field.TypeTime, Nullable: true},
		{Name: "is_blocked", Type: field.TypeBool, Default: false},
		{Name: "is_service", Type: field.TypeBool, Default: false},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "chat_messages", Type: field.TypeUUID, Nullable: true},
		{Name: "problem_messages", Type: field.TypeUUID, Nullable: true},
	}
	// MessagesTable holds the schema information for the "messages" table.
	MessagesTable = &schema.Table{
		Name:       "messages",
		Columns:    MessagesColumns,
		PrimaryKey: []*schema.Column{MessagesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "messages_chats_messages",
				Columns:    []*schema.Column{MessagesColumns[9]},
				RefColumns: []*schema.Column{ChatsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:     "messages_problems_messages",
				Columns:    []*schema.Column{MessagesColumns[10]},
				RefColumns: []*schema.Column{ProblemsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// ProblemsColumns holds the columns for the "problems" table.
	ProblemsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "manager_id", Type: field.TypeUUID},
		{Name: "resolved_at", Type: field.TypeTime, Nullable: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "chat_problems", Type: field.TypeUUID},
	}
	// ProblemsTable holds the schema information for the "problems" table.
	ProblemsTable = &schema.Table{
		Name:       "problems",
		Columns:    ProblemsColumns,
		PrimaryKey: []*schema.Column{ProblemsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "problems_chats_problems",
				Columns:    []*schema.Column{ProblemsColumns[4]},
				RefColumns: []*schema.Column{ChatsColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		ChatsTable,
		MessagesTable,
		ProblemsTable,
	}
)

func init() {
	MessagesTable.ForeignKeys[0].RefTable = ChatsTable
	MessagesTable.ForeignKeys[1].RefTable = ProblemsTable
	ProblemsTable.ForeignKeys[0].RefTable = ChatsTable
}
