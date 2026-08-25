#include "object.h"
#include <cstring>
#include <memory>
#include "context-macros.h"
#include "deps/include/v8-container.h"
#include "deps/include/v8-date.h"
#include "deps/include/v8-object.h"
#include "deps/include/v8-regexp.h"
#include "deps/include/v8-typed-array.h"
#include "isolate-macros.h"
#include "utils.h"
#include "value-macros.h"
#include "value.h"

using namespace v8;

/********** Object **********/

#define LOCAL_OBJECT(ptr) \
  LOCAL_VALUE(ptr)        \
  Local<Object> obj = value.As<Object>()

void ObjectSet(ValuePtr ptr, const char* key, ValuePtr prop_val) {
  LOCAL_OBJECT(ptr);
  Local<String> key_val =
      String::NewFromUtf8(iso, key, NewStringType::kNormal).ToLocalChecked();
  obj->Set(local_ctx, key_val, prop_val->ptr.Get(iso)).Check();
}

void ObjectSetAnyKey(ValuePtr ptr, ValuePtr key, ValuePtr prop_val) {
  LOCAL_OBJECT(ptr);
  Local<Value> local_key = key->ptr.Get(iso);
  obj->Set(local_ctx, local_key, prop_val->ptr.Get(iso)).Check();
}

void ObjectSetIdx(ValuePtr ptr, uint32_t idx, ValuePtr prop_val) {
  LOCAL_OBJECT(ptr);
  obj->Set(local_ctx, idx, prop_val->ptr.Get(iso)).Check();
}

int ObjectSetInternalField(ValuePtr ptr, int idx, ValuePtr val_ptr) {
  LOCAL_OBJECT(ptr);
  m_value* prop_val = static_cast<m_value*>(val_ptr);

  if (idx >= obj->InternalFieldCount()) {
    return 0;
  }

  obj->SetInternalField(idx, prop_val->ptr.Get(iso));

  return 1;
}

int ObjectInternalFieldCount(ValuePtr ptr) {
  LOCAL_OBJECT(ptr);
  return obj->InternalFieldCount();
}

RtnValue ObjectGet(ValuePtr ptr, const char* key) {
  LOCAL_OBJECT(ptr);
  RtnValue rtn = {};

  Local<String> key_val;
  if (!String::NewFromUtf8(iso, key, NewStringType::kNormal)
           .ToLocal(&key_val)) {
    rtn.error = ExceptionError(try_catch, iso, local_ctx);
    return rtn;
  }
  Local<Value> result;
  if (!obj->Get(local_ctx, key_val).ToLocal(&result)) {
    rtn.error = ExceptionError(try_catch, iso, local_ctx);
    return rtn;
  }
  m_value* new_val = new m_value;
  new_val->id = 0;
  new_val->iso = iso;
  new_val->ctx = ctx;
  new_val->ptr = Global<Value>(iso, result);

  rtn.value = tracked_value(ctx, new_val);
  return rtn;
}

RtnValue ObjectGetAnyKey(ValuePtr ptr, ValuePtr key) {
  LOCAL_OBJECT(ptr);
  RtnValue rtn = {};

  Local<Value> local_key = key->ptr.Get(iso);
  Local<Value> result;
  if (!obj->Get(local_ctx, local_key).ToLocal(&result)) {
    rtn.error = ExceptionError(try_catch, iso, local_ctx);
    return rtn;
  }
  m_value* new_val = new m_value;
  new_val->id = 0;
  new_val->iso = iso;
  new_val->ctx = ctx;
  new_val->ptr = Global<Value>(iso, result);

  rtn.value = tracked_value(ctx, new_val);
  return rtn;
}

RtnValue ObjectGetInternalField(ValuePtr ptr, int idx) {
  LOCAL_OBJECT(ptr);
  RtnValue rtn = {};

  if (idx >= obj->InternalFieldCount()) {
    rtn.error.msg = CopyString("internal field index out of range");
    return rtn;
  }

  Local<Data> result = obj->GetInternalField(idx);
  if (!result->IsValue()) {
    rtn.error.msg = CopyString("the internal field did not contain a Value");
    return rtn;
  }

  m_value* new_val = new m_value;
  new_val->id = 0;
  new_val->iso = iso;
  new_val->ctx = ctx;
  new_val->ptr = Global<Value>(iso, result.As<Value>());

  rtn.value = tracked_value(ctx, new_val);
  return rtn;
}

RtnValue ObjectGetIdx(ValuePtr ptr, uint32_t idx) {
  LOCAL_OBJECT(ptr);
  RtnValue rtn = {};

  Local<Value> result;
  if (!obj->Get(local_ctx, idx).ToLocal(&result)) {
    rtn.error = ExceptionError(try_catch, iso, local_ctx);
    return rtn;
  }
  m_value* new_val = new m_value;
  new_val->id = 0;
  new_val->iso = iso;
  new_val->ctx = ctx;
  new_val->ptr = Global<Value>(iso, result);

  rtn.value = tracked_value(ctx, new_val);
  return rtn;
}

int ObjectHas(ValuePtr ptr, const char* key) {
  LOCAL_OBJECT(ptr);
  Local<String> key_val =
      String::NewFromUtf8(iso, key, NewStringType::kNormal).ToLocalChecked();
  return obj->Has(local_ctx, key_val).ToChecked();
}

int ObjectHasAnyKey(ValuePtr ptr, ValuePtr key) {
  LOCAL_OBJECT(ptr);
  Local<Value> local_key = key->ptr.Get(iso);
  return obj->Has(local_ctx, local_key).ToChecked();
}

int ObjectHasIdx(ValuePtr ptr, uint32_t idx) {
  LOCAL_OBJECT(ptr);
  return obj->Has(local_ctx, idx).ToChecked();
}

int ObjectDelete(ValuePtr ptr, const char* key) {
  LOCAL_OBJECT(ptr);
  Local<String> key_val =
      String::NewFromUtf8(iso, key, NewStringType::kNormal).ToLocalChecked();
  return obj->Delete(local_ctx, key_val).ToChecked();
}

int ObjectDeleteAnyKey(ValuePtr ptr, ValuePtr key) {
  LOCAL_OBJECT(ptr);
  Local<Value> local_key = key->ptr.Get(iso);
  return obj->Delete(local_ctx, local_key).ToChecked();
}

int ObjectDeleteIdx(ValuePtr ptr, uint32_t idx) {
  LOCAL_OBJECT(ptr);
  return obj->Delete(local_ctx, idx).ToChecked();
}

RtnValue ObjectGetPrototype(ValuePtr ptr) {
  LOCAL_OBJECT(ptr);
  RtnValue rtn = {};

  Local<Value> result = obj->GetPrototype();
  m_value* new_val = new m_value;
  new_val->id = 0;
  new_val->iso = iso;
  new_val->ctx = ctx;
  new_val->ptr = Global<Value>(iso, result);

  rtn.value = tracked_value(ctx, new_val);
  return rtn;
}

void ObjectSetPrototype(ValuePtr ptr, ValuePtr proto_ptr) {
  LOCAL_OBJECT(ptr);
  // Local<Context> local_ctx = ctx_ptr->ptr.Get(iso);
  obj->SetPrototype(local_ctx, proto_ptr->ptr.Get(iso)).Check();
}

/********** Array **********/

ValuePtr NewArray(ContextPtr ctx, int length) {
  LOCAL_CONTEXT(ctx);
  Local<Array> arr = Array::New(iso, length);
  m_value* val = new m_value;
  val->id = 0;
  val->iso = iso;
  val->ctx = ctx;
  val->ptr = Global<Value>(iso, arr);
  return tracked_value(ctx, val);
}

int ArrayLength(ValuePtr ptr) {
  LOCAL_VALUE_READONLY(ptr);
  Local<Array> arr = value.As<Array>();
  return arr->Length();
}

RtnValue ObjectGetPropertyNames(ValuePtr ptr) {
  LOCAL_OBJECT(ptr);
  RtnValue rtn = {};
  Local<Array> names;
  if (!obj->GetPropertyNames(local_ctx).ToLocal(&names)) {
    rtn.error = ExceptionError(try_catch, iso, local_ctx);
    return rtn;
  }
  m_value* val = new m_value;
  val->id = 0;
  val->iso = iso;
  val->ctx = ctx;
  val->ptr = Global<Value>(iso, names);
  rtn.value = tracked_value(ctx, val);
  return rtn;
}

RtnValue ObjectGetOwnPropertyNames(ValuePtr ptr) {
  LOCAL_OBJECT(ptr);
  RtnValue rtn = {};
  Local<Array> names;
  if (!obj->GetOwnPropertyNames(local_ctx).ToLocal(&names)) {
    rtn.error = ExceptionError(try_catch, iso, local_ctx);
    return rtn;
  }
  m_value* val = new m_value;
  val->id = 0;
  val->iso = iso;
  val->ctx = ctx;
  val->ptr = Global<Value>(iso, names);
  rtn.value = tracked_value(ctx, val);
  return rtn;
}

int ObjectDefineProperty(ValuePtr ptr,
                         const char* key,
                         ValuePtr val_ptr,
                         int attributes) {
  LOCAL_OBJECT(ptr);
  Local<String> key_val =
      String::NewFromUtf8(iso, key, NewStringType::kNormal).ToLocalChecked();
  Maybe<bool> result = obj->DefineOwnProperty(
      local_ctx, key_val, val_ptr->ptr.Get(iso), (PropertyAttribute)attributes);
  if (result.IsNothing()) {
    return 0;
  }
  return result.ToChecked() ? 1 : 0;
}

ValuePtr NewDate(ContextPtr ctx, double ms) {
  LOCAL_CONTEXT(ctx);
  Local<Value> date;
  if (!Date::New(local_ctx, ms).ToLocal(&date)) {
    return nullptr;
  }
  m_value* val = new m_value;
  val->id = 0;
  val->iso = iso;
  val->ctx = ctx;
  val->ptr = Global<Value>(iso, date);
  return tracked_value(ctx, val);
}

double DateValueOf(ValuePtr ptr) {
  LOCAL_VALUE_READONLY(ptr);
  Local<Date> date = value.As<Date>();
  return date->ValueOf();
}

/********** RegExp **********/

ValuePtr NewRegExp(ContextPtr ctx,
                   const char* pattern,
                   int pattern_len,
                   int flags) {
  LOCAL_CONTEXT(ctx);
  Local<String> pat;
  if (!String::NewFromUtf8(iso, pattern, NewStringType::kNormal, pattern_len)
           .ToLocal(&pat)) {
    return nullptr;
  }
  Local<RegExp> re;
  if (!RegExp::New(local_ctx, pat, static_cast<RegExp::Flags>(flags))
           .ToLocal(&re)) {
    return nullptr;
  }
  m_value* val = new m_value;
  val->id = 0;
  val->iso = iso;
  val->ctx = ctx;
  val->ptr = Global<Value>(iso, re);
  return tracked_value(ctx, val);
}

/********** Map **********/

ValuePtr NewMap(ContextPtr ctx) {
  LOCAL_CONTEXT(ctx);
  Local<Map> map = Map::New(iso);
  m_value* val = new m_value;
  val->id = 0;
  val->iso = iso;
  val->ctx = ctx;
  val->ptr = Global<Value>(iso, map);
  return tracked_value(ctx, val);
}

RtnValue MapGet(ValuePtr ptr, ValuePtr key) {
  LOCAL_VALUE(ptr);
  Local<Map> map = value.As<Map>();
  RtnValue rtn = {};
  Local<Value> result;
  if (!map->Get(local_ctx, key->ptr.Get(iso)).ToLocal(&result)) {
    rtn.error = ExceptionError(try_catch, iso, local_ctx);
    return rtn;
  }
  m_value* new_val = new m_value;
  new_val->id = 0;
  new_val->iso = iso;
  new_val->ctx = ctx;
  new_val->ptr = Global<Value>(iso, result);
  rtn.value = tracked_value(ctx, new_val);
  return rtn;
}

void MapSet(ValuePtr ptr, ValuePtr key, ValuePtr val_ptr) {
  LOCAL_VALUE(ptr);
  Local<Map> map = value.As<Map>();
  map->Set(local_ctx, key->ptr.Get(iso), val_ptr->ptr.Get(iso))
      .ToLocalChecked();
}

int MapHas(ValuePtr ptr, ValuePtr key) {
  LOCAL_VALUE(ptr);
  Local<Map> map = value.As<Map>();
  return map->Has(local_ctx, key->ptr.Get(iso)).ToChecked() ? 1 : 0;
}

int MapDelete(ValuePtr ptr, ValuePtr key) {
  LOCAL_VALUE(ptr);
  Local<Map> map = value.As<Map>();
  return map->Delete(local_ctx, key->ptr.Get(iso)).ToChecked() ? 1 : 0;
}

int MapSize(ValuePtr ptr) {
  LOCAL_VALUE_READONLY(ptr);
  Local<Map> map = value.As<Map>();
  return map->Size();
}

/********** Set **********/

ValuePtr NewSet(ContextPtr ctx) {
  LOCAL_CONTEXT(ctx);
  Local<Set> set = Set::New(iso);
  m_value* val = new m_value;
  val->id = 0;
  val->iso = iso;
  val->ctx = ctx;
  val->ptr = Global<Value>(iso, set);
  return tracked_value(ctx, val);
}

void SetAdd(ValuePtr ptr, ValuePtr val_ptr) {
  LOCAL_VALUE(ptr);
  Local<Set> set = value.As<Set>();
  set->Add(local_ctx, val_ptr->ptr.Get(iso)).ToLocalChecked();
}

int SetHas(ValuePtr ptr, ValuePtr val_ptr) {
  LOCAL_VALUE(ptr);
  Local<Set> set = value.As<Set>();
  return set->Has(local_ctx, val_ptr->ptr.Get(iso)).ToChecked() ? 1 : 0;
}

int SetDelete(ValuePtr ptr, ValuePtr val_ptr) {
  LOCAL_VALUE(ptr);
  Local<Set> set = value.As<Set>();
  return set->Delete(local_ctx, val_ptr->ptr.Get(iso)).ToChecked() ? 1 : 0;
}

int SetSize(ValuePtr ptr) {
  LOCAL_VALUE_READONLY(ptr);
  Local<Set> set = value.As<Set>();
  return set->Size();
}

/********** TypedArray **********/

ValuePtr NewArrayBufferFromBytes(ContextPtr ctx, const void* data, int length) {
  LOCAL_CONTEXT(ctx);
  std::unique_ptr<BackingStore> bs = ArrayBuffer::NewBackingStore(iso, length);
  if (data != nullptr && length > 0) {
    memcpy(bs->Data(), data, length);
  }
  Local<ArrayBuffer> ab = ArrayBuffer::New(iso, std::move(bs));
  m_value* val = new m_value;
  val->id = 0;
  val->iso = iso;
  val->ctx = ctx;
  val->ptr = Global<Value>(iso, ab);
  return tracked_value(ctx, val);
}

ValuePtr NewUint8ArrayFromBytes(ContextPtr ctx, const void* data, int length) {
  LOCAL_CONTEXT(ctx);
  std::unique_ptr<BackingStore> bs = ArrayBuffer::NewBackingStore(iso, length);
  if (data != nullptr && length > 0) {
    memcpy(bs->Data(), data, length);
  }
  Local<ArrayBuffer> ab = ArrayBuffer::New(iso, std::move(bs));
  Local<Uint8Array> arr = Uint8Array::New(ab, 0, length);
  m_value* val = new m_value;
  val->id = 0;
  val->iso = iso;
  val->ctx = ctx;
  val->ptr = Global<Value>(iso, arr);
  return tracked_value(ctx, val);
}
