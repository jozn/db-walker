{{range .Tables }}
{{if .PrimaryKey}}
delimiter |
CREATE TRIGGER {{.TableName}}_Create AFTER INSERT ON {{.TableNameSql}}
  FOR EACH ROW
  BEGIN
    INSERT INTO trigger_log (TableName,ChangeType,TargetId) VALUES ("{{.TableNameSql}}","insert",NEW.{{.PrimaryKey.ColumnName}});
#     DELETE FROM test3 WHERE a3 = NEW.a1;
#     UPDATE test4 SET b4 = b4 + 1 WHERE a4 = NEW.a1;
  END;
|
delimiter ;
{{end}}

{{end}}