-- Migration to create the email_templates table
CREATE TABLE IF NOT EXISTS email_templates (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    subject VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_email_template_name ON email_templates(name);

-- Trigger to auto-update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_email_template_updated_at
BEFORE UPDATE ON email_templates
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-- Insert default template
INSERT INTO email_templates (name, subject, body) VALUES
(
    'VERIFY_EMAIL',
    'Verify your email for HomeMie',
    $BODY$
<!doctype html>
<html lang="en" xmlns="http://www.w3.org/1999/xhtml">
<head>
  <meta charset="utf-8">
  <meta name="x-apple-disable-message-reformatting">
  <meta name="viewport" content="width=device-width,initial-scale=1">
  <meta name="color-scheme" content="light dark">
  <meta name="supported-color-schemes" content="light dark">
  <title>Verify your email</title>
  <style>
    /* Mobile tweaks */
    @media (max-width:600px){
      .container{width:100%!important}
      .px{padding-left:20px!important;padding-right:20px!important}
      .btn{display:block!important;width:100%!important}
    }
    /* Dark mode (supported clients) */
    @media (prefers-color-scheme: dark) {
      body, .bg { background:#0b0f14 !important; color:#e7eef7 !important; }
      .card { background:#121922 !important; border-color:#1e2a39 !important; }
      .muted { color:#9fb3c8 !important; }
      a.btn { background:#2a7fff !important; color:#ffffff !important; }
    }
  </style>
</head>
<body style="margin:0;padding:0;background:#f4f6f8;color:#1f2937;">
  <!-- Preheader -->
  <div style="display:none;max-height:0;overflow:hidden;opacity:0;">
    Confirm your email to finish setting up your account.
  </div>

  <table role="presentation" cellpadding="0" cellspacing="0" border="0" width="100%" class="bg" style="background:#f4f6f8;">
    <tr>
      <td align="center" style="padding:32px 12px;">
        <table role="presentation" cellpadding="0" cellspacing="0" border="0" width="600" class="container" style="width:600px;max-width:600px;">
          <!-- Header -->
          <tr>
            <td align="left" class="px" style="padding:24px 32px;">
              <a href="#" style="text-decoration:none;color:#111827;font:600 18px/1.2 -apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Arial,sans-serif;">HomeMie</a>
            </td>
          </tr>

          <!-- Card -->
          <tr>
            <td class="px" style="padding:0 32px;">
              <table role="presentation" cellpadding="0" cellspacing="0" border="0" width="100%" class="card" style="background:#ffffff;border:1px solid #e5e7eb;border-radius:12px;overflow:hidden;">
                <tr>
                  <td style="padding:32px 32px 8px 32px;">
                    <h1 style="margin:0 0 8px 0;font:700 24px/1.25 -apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Arial,sans-serif;color:#111827;">Verify your email</h1>
                    <p style="margin:0 0 16px 0;font:400 16px/1.6 -apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Arial,sans-serif;color:#374151;">
                      Hi {{ .Name }}, thanks for signing up! Please confirm this email address to activate your account.
                    </p>
                  </td>
                </tr>

                <!-- Button -->
                <tr>
                  <td align="center" style="padding:8px 32px 24px 32px;">
                    <!--[if mso]>
                    <v:roundrect xmlns:v="urn:schemas-microsoft-com:vml" href="{{ .VerifyURL }}" style="height:48px;v-text-anchor:middle;width:260px;" arcsize="12%" stroke="f" fillcolor="#2563eb">
                      <w:anchorlock/>
                      <center style="color:#ffffff;font-family:Segoe UI, Arial, sans-serif;font-size:16px;font-weight:600;">Confirm Email</center>
                    </v:roundrect>
                    <![endif]-->
                    <!--[if !mso]><!-- -->
                    <a href="{{ .VerifyURL }}" class="btn"
                       style="display:inline-block;background:#2563eb;color:#ffffff;text-decoration:none;
                              font:600 16px/48px -apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Arial,sans-serif;
                              padding:0 28px;border-radius:8px;min-width:220px;text-align:center;">
                      Confirm Email
                    </a>
                    <!--<![endif]-->
                  </td>
                </tr>

                <!-- Secondary note -->
                <tr>
                  <td style="padding:0 32px 28px 32px;">
                    <p class="muted" style="margin:0 0 10px 0;font:400 14px/1.6 -apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Arial,sans-serif;color:#6b7280;">
                      This link will expire in 30 minutes. If it expires, you can request a new one from the sign-in page.
                    </p>
                    <p class="muted" style="margin:0;font:400 14px/1.6 -apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Arial,sans-serif;color:#6b7280;word-break:break-all;">
                      Having trouble? Copy and paste this URL into your browser:<br>
                      <a href="{{ .VerifyURL }}" style="color:#2563eb;text-decoration:underline;">{{ .VerifyURL }}</a>
                    </p>
                  </td>
                </tr>
              </table>
            </td>
          </tr>

          <!-- Footer -->
          <tr>
            <td class="px" style="padding:16px 32px 0 32px;">
              <p class="muted" style="margin:0 0 4px 0;font:400 12px/1.6 -apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Arial,sans-serif;color:#6b7280;">
                You’re receiving this email because you created an account at HomeMie.
                If this wasn’t you, please ignore this message or contact us at {{ .SupportEmail }}.
              </p>
              <p class="muted" style="margin:6px 0 0 0;font:400 12px/1.6 -apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Arial,sans-serif;color:#9ca3af;">
                © {{ .Year }} HomeMie Inc., 123 Example St, City
              </p>
            </td>
          </tr>

        </table>
      </td>
    </tr>
  </table>
</body>
</html>
$BODY$
);
