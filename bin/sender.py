import smtplib
from email import encoders
from email.mime.base import MIMEBase
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText
import subprocess, datetime, sys


lastweek = datetime.date.today() - datetime.timedelta(weeks=1)
output_name = "output_{}.json".format(lastweek.strftime("%Y%m%d"))

subprocess.run([
    "./cahl.freebsd",
    "-d",
    "output",
    "-D",
    output_name,
    "-e",
    "cahl.xlsx"
    ])

with open("cahl.xlsx", "rb") as attachment:
    # Add the attachment to the message
    pool_file = MIMEBase("application", "octet-stream")
    pool_file.set_payload(attachment.read())
encoders.encode_base64(pool_file)
pool_file.add_header(
    "Content-Disposition",
    f"attachment; filename=cahl.xlsx",
)

with open(output_name, "r") as attachment:
    # Add the attachment to the message
    output_file = MIMEBase("application", "octet-stream")
    output_file.set_payload(attachment.read())
encoders.encode_base64(output_file)
output_file.add_header(
    "Content-Disposition",
    f"attachment; filename=output.json",
)

def send_email(subject, body, sender_name, sender, recipients, password, attachments):
    msg = MIMEMultipart()
    msg['Subject'] = subject
    msg['From'] = sender_name
    msg['To'] = ', '.join(recipients)

    html_part = MIMEText(body)
    msg.attach(html_part)

    for a in attachments:
        msg.attach(a)

    with smtplib.SMTP_SSL('smtp.gmail.com', 465) as smtp_server:
       smtp_server.login(sender, password)
       smtp_server.sendmail(sender, recipients, msg.as_string())
    print("Message sent!")

subject = "Pool de la semaine"
body = "Voici le pool de la semaine.."
sender_name = "Pool Manager"
sender = "roach4facebook@gmail.com"
password = ""
# TODO ^^ add password

recipients = [""]
# TODO ^^ add dad's email
send_email(subject, body, sender_name, sender, recipients, password, [pool_file])

recipients = [""]
# TODO ^^ add my email
send_email(subject, body, sender_name, sender, recipients, password, [pool_file, output_file])