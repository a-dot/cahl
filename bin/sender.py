import os
import smtplib
from email import encoders
from email.mime.base import MIMEBase
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText
import subprocess, datetime, sys
import syslog

syslog.openlog("cahl", syslog.LOG_PID, syslog.LOG_USER)

syslog.syslog("starting")

lastweek = datetime.date.today() - datetime.timedelta(weeks=1)
output_name = "output_{}.json".format(lastweek.strftime("%Y%m%d"))

rc = subprocess.run([
    "../cmd/cahl/cahl",
    "-d",
    "output",
    "-D",
    output_name,
    "-e",
    "cahl.xlsx",
    "-s",
    "20262027"
    ])

if rc.returncode != 0:
    print("error generating pool")
    sys.exit(1)

syslog.syslog(syslog.LOG_DEBUG, f"done running binary, rc={rc.returncode}")

with open("cahl.xlsx", "rb") as attachment:
    # Add the attachment to the message
    pool_file = MIMEBase("application", "octet-stream")
    pool_file.set_payload(attachment.read())
encoders.encode_base64(pool_file)
pool_file.add_header(
    "Content-Disposition",
    f"attachment; filename=cahl.xlsx",
)

syslog.syslog(syslog.LOG_DEBUG, "done reading excel sheet")

with open(output_name, "r") as attachment:
    # Add the attachment to the message
    output_file = MIMEBase("application", "octet-stream")
    output_file.set_payload(attachment.read())
encoders.encode_base64(output_file)
output_file.add_header(
    "Content-Disposition",
    f"attachment; filename=output.json",
)

syslog.syslog(syslog.LOG_DEBUG, "done reading json file")

def send_email(subject, body, sender_name, sender, recipients, password, attachments):
    msg = MIMEMultipart()
    msg['Subject'] = subject
    msg['From'] = sender_name
    msg['To'] = ', '.join(recipients)

    html_part = MIMEText(body)
    msg.attach(html_part)

    for a in attachments:
        msg.attach(a)

    syslog.syslog(syslog.LOG_DEBUG, f"sending message to={recipients}, to_formatted={msg['To']}")

    with smtplib.SMTP_SSL('smtp.gmail.com', 465) as smtp_server:
       smtp_server.login(sender, password)
       syslog.syslog(syslog.LOG_DEBUG, "logged in")
       
       smtp_server.sendmail(sender, recipients, msg.as_string())
       syslog.syslog(syslog.LOG_DEBUG, "message sent")

       smtp_server.quit()

    print("Message sent!")

subject = "Pool de la semaine"
body = "Voici le pool de la semaine.."

sender_name = os.environ.get("SMTP_SENDER_NAME")
if not sender_name:
    print("SMTP_SENDER_NAME environment variable not set")
    sys.exit(1)

sender = os.environ.get("SMTP_SENDER")
if not sender:
    print("SMTP_SENDER environment variable not set")
    sys.exit(1)

password = os.environ.get("SMTP_PASS")
if not password:
    print("SMTP_PASS environment variable not set")
    sys.exit(1)

syslog.syslog(syslog.LOG_DEBUG, "send first email")
recipients = os.environ.get("SMTP_RECIPIENT")
if not recipients:
    print("SMTP_RECIPIENT environment variable not set")
    sys.exit(1)

send_email(subject, body, sender_name, sender, recipients, password, [pool_file])

syslog.syslog(syslog.LOG_DEBUG, "send second email")
recipients = os.environ.get("SMTP_RECIPIENT_DEV")
if not recipients:
    print("SMTP_RECIPIENT_DEV environment variable not set")
    sys.exit(1)
send_email(subject, body, sender_name, sender, recipients, password, [pool_file, output_file])

syslog.closelog