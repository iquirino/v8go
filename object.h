#ifndef V8GO_OBJECT_H
#define V8GO_OBJECT_H

#include <stdint.h>

#include "errors.h"

#ifdef __cplusplus

extern "C" {
#endif

typedef struct m_ctx m_ctx;
typedef m_ctx* ContextPtr;

extern void ObjectSet(ValuePtr ptr, const char* key, ValuePtr val_ptr);
extern void ObjectSetAnyKey(ValuePtr ptr, ValuePtr key, ValuePtr val_ptr);
extern void ObjectSetIdx(ValuePtr ptr, uint32_t idx, ValuePtr val_ptr);
extern int ObjectSetInternalField(ValuePtr ptr, int idx, ValuePtr val_ptr);
extern int ObjectInternalFieldCount(ValuePtr ptr);

extern RtnValue ObjectGet(ValuePtr ptr, const char* key);
extern RtnValue ObjectGetAnyKey(ValuePtr ptr, ValuePtr key);
extern RtnValue ObjectGetIdx(ValuePtr ptr, uint32_t idx);
extern RtnValue ObjectGetInternalField(ValuePtr ptr, int idx);
int ObjectHas(ValuePtr ptr, const char* key);
int ObjectHasAnyKey(ValuePtr ptr, ValuePtr key);
int ObjectHasIdx(ValuePtr ptr, uint32_t idx);
int ObjectDelete(ValuePtr ptr, const char* key);
int ObjectDeleteAnyKey(ValuePtr ptr, ValuePtr key);
int ObjectDeleteIdx(ValuePtr ptr, uint32_t idx);
extern RtnValue ObjectGetPrototype(ValuePtr ptr);
extern void ObjectSetPrototype(ValuePtr ptr, ValuePtr proto_ptr);

extern ValuePtr NewArray(ContextPtr ctx, int length);
extern int ArrayLength(ValuePtr ptr);

extern RtnValue ObjectGetPropertyNames(ValuePtr ptr);
extern RtnValue ObjectGetOwnPropertyNames(ValuePtr ptr);
extern int ObjectDefineProperty(ValuePtr ptr, const char* key, ValuePtr val_ptr, int attributes);
extern ValuePtr NewDate(ContextPtr ctx, double ms);
extern double DateValueOf(ValuePtr ptr);
extern ValuePtr NewRegExp(ContextPtr ctx, const char* pattern, int pattern_len, int flags);

extern ValuePtr NewMap(ContextPtr ctx);
extern RtnValue MapGet(ValuePtr ptr, ValuePtr key);
extern void MapSet(ValuePtr ptr, ValuePtr key, ValuePtr val_ptr);
extern int MapHas(ValuePtr ptr, ValuePtr key);
extern int MapDelete(ValuePtr ptr, ValuePtr key);
extern int MapSize(ValuePtr ptr);

extern ValuePtr NewSet(ContextPtr ctx);
extern void SetAdd(ValuePtr ptr, ValuePtr val_ptr);
extern int SetHas(ValuePtr ptr, ValuePtr val_ptr);
extern int SetDelete(ValuePtr ptr, ValuePtr val_ptr);
extern int SetSize(ValuePtr ptr);

extern ValuePtr NewArrayBufferFromBytes(ContextPtr ctx, const void* data, int length);
extern ValuePtr NewUint8ArrayFromBytes(ContextPtr ctx, const void* data, int length);

#ifdef __cplusplus
}  // extern "C"
#endif
#endif
