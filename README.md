<!-- togo-brand -->
<p align="center">
  <img src=".github/assets/togo-mark.svg" width="96" alt="togo" />
</p>
<h1 align="center">mail</h1>
<p align="center"><sub>part of the <a href="https://github.com/togo-framework">togo-framework</a> — the full-stack Go + React framework</sub></p>

Email for [togo](https://github.com/togo-framework/togo): a Mailer contract with an
SMTP driver (default) and a dev log driver. AWS SES, Resend, etc. ship as driver
plugins that call mail.RegisterDriver.

```bash
togo install togo-framework/mail
```

Env: MAIL_DRIVER (smtp|log), MAIL_HOST, MAIL_PORT, MAIL_USERNAME, MAIL_PASSWORD,
MAIL_FROM, MAIL_ENCRYPTION (tls|starttls|none).


---

## 💎 Premium sponsors

togo is proudly sponsored by **ID8 Media** and **One Studio**.

<p align="center">
  <a href="https://id8media.com"><img src=".github/assets/id8media.svg" height="44" alt="ID8 Media" /></a>
  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
  <a href="https://one-studio.co"><img src=".github/assets/one-studio.jpeg" height="44" alt="One Studio" /></a>
</p>
