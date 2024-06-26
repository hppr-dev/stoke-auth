# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc
import warnings

from proto import cargo_pb2 as proto_dot_cargo__pb2

GRPC_GENERATED_VERSION = '1.64.1'
GRPC_VERSION = grpc.__version__
EXPECTED_ERROR_RELEASE = '1.65.0'
SCHEDULED_RELEASE_DATE = 'June 25, 2024'
_version_not_supported = False

try:
    from grpc._utilities import first_version_is_lower
    _version_not_supported = first_version_is_lower(GRPC_VERSION, GRPC_GENERATED_VERSION)
except ImportError:
    _version_not_supported = True

if _version_not_supported:
    warnings.warn(
        f'The grpc package installed is at version {GRPC_VERSION},'
        + f' but the generated code in proto/cargo_pb2_grpc.py depends on'
        + f' grpcio>={GRPC_GENERATED_VERSION}.'
        + f' Please upgrade your grpc module to grpcio>={GRPC_GENERATED_VERSION}'
        + f' or downgrade your generated code using grpcio-tools<={GRPC_VERSION}.'
        + f' This warning will become an error in {EXPECTED_ERROR_RELEASE},'
        + f' scheduled for release on {SCHEDULED_RELEASE_DATE}.',
        RuntimeWarning
    )


class CargoHoldStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.GetContents = channel.unary_unary(
                '/cargo.CargoHold/GetContents',
                request_serializer=proto_dot_cargo__pb2.ContentRequest.SerializeToString,
                response_deserializer=proto_dot_cargo__pb2.ContentReply.FromString,
                _registered_method=True)
        self.LoadItems = channel.stream_stream(
                '/cargo.CargoHold/LoadItems',
                request_serializer=proto_dot_cargo__pb2.LoadItemRequest.SerializeToString,
                response_deserializer=proto_dot_cargo__pb2.LoadItemReply.FromString,
                _registered_method=True)


class CargoHoldServicer(object):
    """Missing associated documentation comment in .proto file."""

    def GetContents(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def LoadItems(self, request_iterator, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_CargoHoldServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'GetContents': grpc.unary_unary_rpc_method_handler(
                    servicer.GetContents,
                    request_deserializer=proto_dot_cargo__pb2.ContentRequest.FromString,
                    response_serializer=proto_dot_cargo__pb2.ContentReply.SerializeToString,
            ),
            'LoadItems': grpc.stream_stream_rpc_method_handler(
                    servicer.LoadItems,
                    request_deserializer=proto_dot_cargo__pb2.LoadItemRequest.FromString,
                    response_serializer=proto_dot_cargo__pb2.LoadItemReply.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'cargo.CargoHold', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))
    server.add_registered_method_handlers('cargo.CargoHold', rpc_method_handlers)


 # This class is part of an EXPERIMENTAL API.
class CargoHold(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def GetContents(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/cargo.CargoHold/GetContents',
            proto_dot_cargo__pb2.ContentRequest.SerializeToString,
            proto_dot_cargo__pb2.ContentReply.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def LoadItems(request_iterator,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.stream_stream(
            request_iterator,
            target,
            '/cargo.CargoHold/LoadItems',
            proto_dot_cargo__pb2.LoadItemRequest.SerializeToString,
            proto_dot_cargo__pb2.LoadItemReply.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)
