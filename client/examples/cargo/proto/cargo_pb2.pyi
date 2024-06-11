from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Item(_message.Message):
    __slots__ = ("name", "id", "width", "height", "depth")
    NAME_FIELD_NUMBER: _ClassVar[int]
    ID_FIELD_NUMBER: _ClassVar[int]
    WIDTH_FIELD_NUMBER: _ClassVar[int]
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    DEPTH_FIELD_NUMBER: _ClassVar[int]
    name: str
    id: str
    width: int
    height: int
    depth: int
    def __init__(self, name: _Optional[str] = ..., id: _Optional[str] = ..., width: _Optional[int] = ..., height: _Optional[int] = ..., depth: _Optional[int] = ...) -> None: ...

class ContentRequest(_message.Message):
    __slots__ = ("content_filter",)
    CONTENT_FILTER_FIELD_NUMBER: _ClassVar[int]
    content_filter: str
    def __init__(self, content_filter: _Optional[str] = ...) -> None: ...

class ContentReply(_message.Message):
    __slots__ = ("filter_matched", "items")
    FILTER_MATCHED_FIELD_NUMBER: _ClassVar[int]
    ITEMS_FIELD_NUMBER: _ClassVar[int]
    filter_matched: bool
    items: _containers.RepeatedCompositeFieldContainer[Item]
    def __init__(self, filter_matched: bool = ..., items: _Optional[_Iterable[_Union[Item, _Mapping]]] = ...) -> None: ...

class LoadItemRequest(_message.Message):
    __slots__ = ("num_items", "item")
    NUM_ITEMS_FIELD_NUMBER: _ClassVar[int]
    ITEM_FIELD_NUMBER: _ClassVar[int]
    num_items: int
    item: Item
    def __init__(self, num_items: _Optional[int] = ..., item: _Optional[_Union[Item, _Mapping]] = ...) -> None: ...

class LoadItemReply(_message.Message):
    __slots__ = ("loaded", "message")
    LOADED_FIELD_NUMBER: _ClassVar[int]
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    loaded: bool
    message: str
    def __init__(self, loaded: bool = ..., message: _Optional[str] = ...) -> None: ...
