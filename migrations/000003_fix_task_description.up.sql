DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'todoapp' AND table_name = 'tasks' AND column_name = 'desciption'
    ) AND NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'todoapp' AND table_name = 'tasks' AND column_name = 'description'
    ) THEN
        ALTER TABLE todoapp.tasks RENAME COLUMN desciption TO description;
    END IF;
END $$;
