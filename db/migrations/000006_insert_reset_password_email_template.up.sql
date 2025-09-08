INSERT INTO "email_templates" ("name", "subject", "body") VALUES
(
    'RESET_PASSWORD',
    'Reset Your Password',
    'Hi {{.Name}},<br><br>You recently requested to reset your password for your Homemie account. Click the link below to reset it.<br><br><a href="{{.ResetURL}}">Reset Your Password</a><br><br>If you did not request a password reset, please ignore this email or contact support if you have questions.<br><br>Thanks,<br>The Homemie Team'
);
