#include "symbol.h"
#include "context.h"
#include "deps/include/v8-primitive.h"
#include "errors.h"
#include "isolate-macros.h"
#include "utils.h"
#include "value-macros.h"

using namespace v8;

ValuePtr BuiltinSymbol(IsolatePtr iso, SymbolIndex idx) {
  ISOLATE_SCOPE(iso);
  INTERNAL_CONTEXT(iso);
  Local<Symbol> sym;
  switch (idx) {
    case SYMBOL_ASYNC_ITERATOR:
      sym = Symbol::GetAsyncIterator(iso);
      break;
    case SYMBOL_HAS_INSTANCE:
      sym = Symbol::GetHasInstance(iso);
      break;
    case SYMBOL_IS_CONCAT_SPREADABLE:
      sym = Symbol::GetIsConcatSpreadable(iso);
      break;
    case SYMBOL_ITERATOR:
      sym = Symbol::GetIterator(iso);
      break;
    case SYMBOL_MATCH:
      sym = Symbol::GetMatch(iso);
      break;
    case SYMBOL_REPLACE:
      sym = Symbol::GetReplace(iso);
      break;
    case SYMBOL_SEARCH:
      sym = Symbol::GetSearch(iso);
      break;
    case SYMBOL_SPLIT:
      sym = Symbol::GetSplit(iso);
      break;
    case SYMBOL_TO_PRIMITIVE:
      sym = Symbol::GetToPrimitive(iso);
      break;
    case SYMBOL_TO_STRING_TAG:
      sym = Symbol::GetToStringTag(iso);
      break;
    case SYMBOL_UNSCOPABLES:
      sym = Symbol::GetUnscopables(iso);
      break;
    default:
      return nullptr;
  }
  m_value* val = new m_value;
  val->id = 0;
  val->iso = iso;
  val->ctx = ctx;
  val->ptr = Global<Value>(iso, sym);
  return tracked_value(ctx, val);
}

const char* SymbolDescription(ValuePtr ptr) {
  LOCAL_VALUE(ptr);
  Local<Symbol> sym = value.As<Symbol>();
  Local<Value> descr = sym->Description(iso);
  String::Utf8Value utf8(iso, descr);
  return CopyString(utf8);
}

/********** Private **********/

int ObjectSetPrivate(ValuePtr ptr, const char* key, int key_len, ValuePtr val_ptr) {
  LOCAL_VALUE(ptr);
  Local<Object> obj = value.As<Object>();
  Local<String> str = String::NewFromUtf8(iso, key, NewStringType::kNormal, key_len).ToLocalChecked();
  Local<Private> priv = Private::ForApi(iso, str);
  Maybe<bool> result = obj->SetPrivate(local_ctx, priv, val_ptr->ptr.Get(iso));
  if (result.IsNothing()) return 0;
  return result.ToChecked() ? 1 : 0;
}

RtnValue ObjectGetPrivate(ValuePtr ptr, const char* key, int key_len) {
  LOCAL_VALUE(ptr);
  Local<Object> obj = value.As<Object>();
  RtnValue rtn = {};
  Local<String> str = String::NewFromUtf8(iso, key, NewStringType::kNormal, key_len).ToLocalChecked();
  Local<Private> priv = Private::ForApi(iso, str);
  Local<Value> result;
  if (!obj->GetPrivate(local_ctx, priv).ToLocal(&result)) {
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

int ObjectHasPrivate(ValuePtr ptr, const char* key, int key_len) {
  LOCAL_VALUE(ptr);
  Local<Object> obj = value.As<Object>();
  Local<String> str = String::NewFromUtf8(iso, key, NewStringType::kNormal, key_len).ToLocalChecked();
  Local<Private> priv = Private::ForApi(iso, str);
  Maybe<bool> result = obj->HasPrivate(local_ctx, priv);
  if (result.IsNothing()) return 0;
  return result.ToChecked() ? 1 : 0;
}

int ObjectDeletePrivate(ValuePtr ptr, const char* key, int key_len) {
  LOCAL_VALUE(ptr);
  Local<Object> obj = value.As<Object>();
  Local<String> str = String::NewFromUtf8(iso, key, NewStringType::kNormal, key_len).ToLocalChecked();
  Local<Private> priv = Private::ForApi(iso, str);
  Maybe<bool> result = obj->DeletePrivate(local_ctx, priv);
  if (result.IsNothing()) return 0;
  return result.ToChecked() ? 1 : 0;
}
