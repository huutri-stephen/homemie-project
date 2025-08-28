DROP TRIGGER IF EXISTS trg_update_email_template_updated_at ON email_templates;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP INDEX IF EXISTS idx_email_template_name;
DROP TABLE IF EXISTS email_templates;
