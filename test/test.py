import os


# apis = os.getenv("API_KEY")
aws_key = os.getenv("AWS_SECRET_ID")

domain = os.getenv("DOMAIN")

if all in [aws_key, domain]:
    print(f"we shit the bed nothing is there  \n{aws_key}\n {domain}")

print(f"Secrets Recieved  \n{aws_key}\n {domain}")
