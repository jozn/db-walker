

delimiter |
CREATE TRIGGER chat_Create AFTER INSERT ON chat
  FOR EACH ROW
  BEGIN
    INSERT INTO trigger_log (TableName,ChangeType,TargetId) VALUES ("chat","insert",NEW.ChatKey);
#     DELETE FROM test3 WHERE a3 = NEW.a1;
#     UPDATE test4 SET b4 = b4 + 1 WHERE a4 = NEW.a1;
  END;
|
delimiter ;




delimiter |
CREATE TRIGGER comments_Create AFTER INSERT ON comments
  FOR EACH ROW
  BEGIN
    INSERT INTO trigger_log (TableName,ChangeType,TargetId) VALUES ("comments","insert",NEW.Id);
#     DELETE FROM test3 WHERE a3 = NEW.a1;
#     UPDATE test4 SET b4 = b4 + 1 WHERE a4 = NEW.a1;
  END;
|
delimiter ;




delimiter |
CREATE TRIGGER trigger_log_Create AFTER INSERT ON trigger_log
  FOR EACH ROW
  BEGIN
    INSERT INTO trigger_log (TableName,ChangeType,TargetId) VALUES ("trigger_log","insert",NEW.Id);
#     DELETE FROM test3 WHERE a3 = NEW.a1;
#     UPDATE test4 SET b4 = b4 + 1 WHERE a4 = NEW.a1;
  END;
|
delimiter ;


