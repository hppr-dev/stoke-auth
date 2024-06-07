from grpc import access_token_call_credentials


class GRPCClientCredentials:
    def __init__(self, token : str, refresh_token : str = "", refresh_url : str = ""):
        self.token = token
        self.refresh = refresh_token
        self.refresh_url = refresh_url

    def call_credentials(self):
        return access_token_call_credentials(self.token)
