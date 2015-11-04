package main

/*
#include <stdio.h>
#include <aj_debug.h>
#include <aj_guid.h>
#include <aj_creds.h>
#include <aj_nvram.h>
#include "alljoyn.h"
//#include "ServicesCommon.c"

#include "alljoyn/services_common/PropertyStore.h"
#include "NotificationCommon/NotificationCommon.c"
#include "NotificationProducer/NotificationProducer.c"

AJ_BusAttachment c_bus;
AJ_Message c_message;
AJ_Message c_reply;
void * c_propGetter;
AJ_Arg c_arg;
AJ_SessionOpts session_opts = { AJ_SESSION_TRAFFIC_MESSAGES, AJ_SESSION_PROXIMITY_ANY, AJ_TRANSPORT_ANY, TRUE };

// from "ServicesCommon.c"
AJ_Status AJSVC_MarshalAppId(AJ_Message* msg, const char* appId)
{
    AJ_Status status;
    uint8_t binAppId[UUID_LENGTH];
    uint32_t sz = strlen(appId);

    if (sz > UUID_LENGTH * 2) { // Crop application id that is too long
        sz = UUID_LENGTH * 2;
    }
    status = AJ_HexToRaw(appId, sz, binAppId, UUID_LENGTH);
    if (status != AJ_OK) {
        return status;
    }
    status = AJ_MarshalArgs(msg, APP_ID_SIGNATURE, binAppId, sz / 2);

    return status;
}

int8_t AJSVC_PropertyStore_GetLanguageIndex(const char* const language){
	return 0;
}

const char* AJSVC_PropertyStore_GetValueForLang(int8_t fieldIndex, int8_t langIndex){

	 if (fieldIndex <= AJSVC_PROPERTY_STORE_ERROR_FIELD_INDEX || fieldIndex >= AJSVC_PROPERTY_STORE_NUMBER_OF_KEYS) {
        return NULL;
    }

	switch(fieldIndex)
	{
	    case AJSVC_PROPERTY_STORE_DEVICE_ID:
	        return (char*)GetAboutProperty("DeviceId", "en");
	    case AJSVC_PROPERTY_STORE_APP_ID:
	        return "f85bda3742d04ff782c774f01b458cba";
		case AJSVC_PROPERTY_STORE_DEVICE_NAME:
	        return (char*)GetAboutProperty("DeviceIName", "en");
		case AJSVC_PROPERTY_STORE_APP_NAME:
			return (char*)GetAboutProperty("AppName", "en");
	    default :
	       return NULL;
	}
}


#define NUM_TEXTS   1
AJNS_DictionaryEntry textToSend[NUM_TEXTS];

static AJNS_NotificationContent notificationContent;

void SendNotification(uint16_t messageType, char * lang, char * msg)
{
//    uint16_t messageType = AJNS_NOTIFICATION_MESSAGE_TYPE_INFO;
    uint32_t ttl = AJNS_NOTIFICATION_TTL_MIN; // Note needs to be in the range AJNS_NOTIFICATION_TTL_MIN..AJNS_NOTIFICATION_TTL_MAX
    uint32_t serialNum;

    notificationContent.numTexts = NUM_TEXTS;
    textToSend[0].key   = lang;
    textToSend[0].value = msg;
    notificationContent.texts = textToSend;

    AJNS_Producer_SendNotification(&c_bus, &notificationContent, messageType, ttl, &serialNum);
}


const char* AJSVC_PropertyStore_GetValue(int8_t fieldIndex)
{
    return AJSVC_PropertyStore_GetValueForLang(fieldIndex, AJSVC_PROPERTY_STORE_NO_LANGUAGE_INDEX);
}

void * Get_Session_Opts() {
	return &session_opts;
}

void * Get_Arg() {
	return &c_arg;
}

AJ_Status AJ_MarshalArgs_cgo(AJ_Message* msg, char * a, char * b, char * c, char * d) {
	return AJ_MarshalArgs(msg, a, b, c, d);
}

AJ_Status MarshalArg(AJ_Message* msg, char * sig, void * value) {
	printf("SIG: %s\n", sig);
	return AJ_MarshalArgs(msg, sig, value);
}

AJ_Status AJ_MarshalSignal_cgo(AJ_Message* msg, uint32_t msgId, uint32_t sessionId, uint8_t flags, uint32_t ttl) {
	return AJ_MarshalSignal(&c_bus, msg, msgId, NULL, (AJ_SessionId) sessionId, flags, ttl);
}


uint32_t Get_AJ_Message_msgId() {
	return c_message.msgId;
}

uint32_t Get_AJ_Message_bodyLen() {
	return c_message.hdr->bodyLen;
}

const char * Get_AJ_Message_signature() {
	return c_message.signature;
}

const char * Get_AJ_Message_objPath() {
   return c_message.objPath;
}

const char * Get_AJ_Message_iface() {
   return c_message.iface;
}

const char * Get_AJ_Message_member() {
   return c_message.member;
}

const char * Get_AJ_Message_destination() {
   return c_message.destination;
}

AJ_Message * Get_AJ_ReplyMessage() {
	return &c_reply;
}

AJ_Message * Get_AJ_Message() {
	return &c_message;
}
AJ_BusAttachment * Get_AJ_BusAttachment() {
	return &c_bus;
}

AJ_Object * Allocate_AJ_Object_Array(uint32_t array_size) {
	return (AJ_Object*)AJ_Malloc(sizeof(AJ_Object)*array_size);
}

void * Create_AJ_Object(uint32_t index, AJ_Object * array, char* path, AJ_InterfaceDescription* interfaces, uint8_t flags, void* context) {
	// BE CAREFULL WHEN YOU WILL IMPLEMENT OBJECT DELETION
	// YOU MUST DELETE ALL ALLOCS
	AJ_Object * obj = array + index * sizeof(AJ_Object);
	if(path) {
		char *c = AJ_Malloc(strlen(path) + 1);
		strcpy(c, path);
		obj->path = c;
	} else {
		obj->path = 0;
	}

	if(interfaces) {
		int ic = 0;
		while(interfaces[ic++]);
		AJ_InterfaceDescription *interfacescopy = AJ_Malloc(ic * sizeof(AJ_InterfaceDescription*));
		int i;
		for(i = 0; i < ic; i++) {
			if(interfaces[i]) {
				int iic = 0;
				while(interfaces[i][iic++]);
				char **newitem = AJ_Malloc(iic * sizeof(char *));
				int j;
				for(j = 0; j < iic; j++) {
					if(interfaces[i][j]) {
						char *c = AJ_Malloc(strlen(interfaces[i][j]) + 1);
						strcpy(c, interfaces[i][j]);
						newitem[j] = c;
					} else {
						newitem[j] = 0;
					}
				}
				interfacescopy[i] = (AJ_InterfaceDescription)newitem;
			} else {
				interfacescopy[i] = 0;
			}
		}
		obj->interfaces = interfacescopy;
	} else {
		obj->interfaces = 0;
	}

	obj->flags = flags;
	obj->context = context;
	return obj;
}

AJ_Status MyAboutPropGetter_cgo(AJ_Message* reply, const char* language) {
	//printf("C.MyAboutPropGetter_cgo() called\n");
	return MyAboutPropGetter(reply, language);
}

AJ_Status UnmarshalJoinSessionArgs(AJ_Message* msg, uint16_t * port, uint32_t * sessionId) {
	char* joiner;
	return AJ_UnmarshalArgs(msg, "qus", port, sessionId, &joiner);
}

AJ_Status UnmarshalLostSessionArgs(AJ_Message* msg, uint32_t * sessionId, uint32_t * reason) {
	return AJ_UnmarshalArgs(msg, "uu", sessionId, reason);
}

int UnmarshalPort() {
	uint16_t port;
	char* joiner;
	uint32_t sessionId;

	AJ_UnmarshalArgs(&c_message, "qus", &port, &sessionId, &joiner);
	return port;
}

*/
import "C"
