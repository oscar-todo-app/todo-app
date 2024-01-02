This folders contains all the files related to k8s.

- After applying the terraform files from the infra folder we can ran the Make and this will create all the necessary k8s definitions.
- Make sure to change the service accounts and the container repo to match the ones in your account:
    - Cert-manager/values.
    - Todo-template/values && secret file, && and service account file.
    - Ingress and dns values to be the ones that match your files.


