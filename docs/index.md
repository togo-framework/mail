# mail

Email for [togo](https://github.com/togo-framework/togo): a Mailer contract with an
SMTP driver (default) and a dev log driver. AWS SES, Resend, etc. ship as driver
plugins that call mail.RegisterDriver.

```bash
togo install togo-framework/mail
```

Env: MAIL_DRIVER (smtp|log), MAIL_HOST, MAIL_PORT, MAIL_USERNAME, MAIL_PASSWORD,
MAIL_FROM, MAIL_ENCRYPTION (tls|starttls|none).
