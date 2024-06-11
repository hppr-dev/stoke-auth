from typing import Any, Callable, Dict, Optional, List
from grpc import HandlerCallDetails, RpcMethodHandler, ServerInterceptor, StatusCode, stream_stream_rpc_method_handler, stream_unary_rpc_method_handler, unary_stream_rpc_method_handler, unary_unary_rpc_method_handler
from stoke.client import StokeClient

UNARY_UNARY   = "unary-unary"
UNARY_STREAM  = "unary-stream"
STREAM_UNARY  = "stream-unary"
STREAM_STREAM = "stream-stream"

class ServerTokenInterceptor(ServerInterceptor):

    def __init__(self, stoke_client : StokeClient, required_claims : Dict[str, str] = {}, endpoint_type : str = UNARY_UNARY, method : str = ""):
        self.client = stoke_client
        self.endpoint_type = endpoint_type
        self.method = method
        self.required_claims = required_claims

    def intercept_service(self, continuation: Callable[[HandlerCallDetails], RpcMethodHandler], handler_call_details: HandlerCallDetails) -> RpcMethodHandler:
        if self.method != "" and self.method == handler_call_details:
            return continuation(handler_call_details)

        claims = self._verify_details(handler_call_details)
        if claims is None or any([claims[k] != self.required_claims[k] for k in self.required_claims]):
            return self._reject_handler()

        return continuation(handler_call_details)

    def _verify_details(self, handler_call_details : HandlerCallDetails) -> Optional[Dict[str, Any]]:
        metadata = handler_call_details.invocation_metadata
        if "authorization" not in metadata:
            return None

        return self.client.parse_token(metadata["authorization"].removeprefix("Bearer "))

    def _reject_handler(self) -> RpcMethodHandler:
        def reject(_, context):
            context.abort(StatusCode.UNAUTHENTICATED, "Unable to authenticate request")

        if self.endpoint_type == UNARY_STREAM:
            return unary_stream_rpc_method_handler(reject)
        
        if self.endpoint_type == STREAM_UNARY:
            return stream_unary_rpc_method_handler(reject)

        if self.endpoint_type == STREAM_STREAM:
            return stream_stream_rpc_method_handler(reject)

        return unary_unary_rpc_method_handler(reject)

def intercept_all(stoke_client : StokeClient, required_claims : Dict[str, str] = {}) -> List[ServerTokenInterceptor]:
    return [ ServerTokenInterceptor(stoke_client, required_claims, et) for et in [UNARY_UNARY, UNARY_STREAM, STREAM_UNARY, STREAM_STREAM] ]
