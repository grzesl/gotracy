#define TRACY_ENABLE

#include "tracy/TracyC.h"
#include "gotracy.h"
#include <map>
#include <string>
#include <stdio.h>

typedef struct  ___tracy_source_location_data TracyCZoneLocation;

struct TracyZoneData
{
    TracyCZoneLocation loc;
    TracyCZoneCtx ctx;  
};

std::map<std::string, TracyZoneData*> TracyCZoneCtxMap;
int TracyCZoneCtxCounter = 0;

TracyZoneData* GetZoneContext(std::string c)
{
    auto search = TracyCZoneCtxMap.find(c);

    if(search == TracyCZoneCtxMap.end())
    {
        auto data = new TracyZoneData();
        TracyCZoneCtxMap[c] = data;
        return data;
    } 
    else {
        return search->second;
    }
   
}

void SetZoneContext(std::string c, TracyZoneData *data)
{
    TracyCZoneCtxMap[c] = data;
}

void DelZoneContext(std::string c)
{
    auto it = TracyCZoneCtxMap.find(c);
    if(it!=TracyCZoneCtxMap.end())
        TracyCZoneCtxMap.erase(it);
}

int IsZoneContextExist(std::string c)
{
    auto search = TracyCZoneCtxMap.find(c);

    if(search == TracyCZoneCtxMap.end())
        return 0;
    return 1;
}

void GoTracySetThreadName(const char*name)
{
    ___tracy_set_thread_name(name);
}

int GoTracyZoneBegin(const char*name,const char *function,const char*file, uint32_t line, uint32_t color)
{
    TracyCZoneCtxCounter++;
    std::string s = std::to_string(TracyCZoneCtxCounter);
    TracyZoneData *data = GetZoneContext(s);
    data->ctx = TracyCZoneCtx {};
    data->loc = TracyCZoneLocation {};
    data->loc.name = name;
    data->loc.function = function;
    data->loc.file = file;
    data->loc.line = line;
    data->loc.color = color;

    data->ctx = ___tracy_emit_zone_begin( (___tracy_source_location_data*)&data->loc, 1);

    char*dbg = new char [128];
    sprintf(dbg, "GoTracyZoneBegin c: %d id: %d\n", TracyCZoneCtxCounter, data->ctx.id);
    GoTracyMessageL(dbg);
    return TracyCZoneCtxCounter;
}

void GoTracyZoneEnd(const char *c){

    char*dbg2 = new char [128];
    sprintf(dbg2, "GoTracyZoneEnd1 c: %s\n", c);
    GoTracyMessageL(dbg2);

    std::string czone = std::string(c);
    if(!IsZoneContextExist(czone))
        {
            char*dbg3 = new char [128];
            sprintf(dbg3, "GoTracyZoneEnd3 not exist c: %s\n",  c);
            GoTracyMessageL(dbg3);
        }
    TracyZoneData *data = GetZoneContext(czone);

    char*dbg = new char [128];
    sprintf(dbg, "GoTracyZoneEnd2 c: %s id: %d\n",  c, data->ctx.id);
    GoTracyMessageL(dbg);

    ___tracy_emit_zone_end(data->ctx);

    DelZoneContext(c);
}

void GoTracyZoneValue(const char* c, uint64_t value){
    TracyZoneData *data = GetZoneContext(std::string(c));
    ___tracy_emit_zone_value(data->ctx, value);
}

void GoTracyZoneText(const char* c, char* text){
    TracyZoneData *data = GetZoneContext(std::string(c));
    ___tracy_emit_zone_text(data->ctx, text, strlen(text));
}

void  GoTracyMessageL(char * msg)
{
    ___tracy_emit_messageL( msg, 1 );
}

void  GoTracyMessageLC(char * msg, uint32_t color)
{
    ___tracy_emit_messageLC( msg, color, 1 );
}

void  GoTracyFrameMark()
{
    ___tracy_emit_frame_mark((char*)0);
}

void  GoTracyFrameMarkName(char *name)
{
    ___tracy_emit_frame_mark(name);
}

void  GoTracyFrameMarkStart(char *name)
{
    ___tracy_emit_frame_mark_start(name);
}

void  GoTracyFrameMarkEnd(char *name)
{
    ___tracy_emit_frame_mark_end(name);
}

void GoTracyPlotFloat(char *name, float val)
{
    ___tracy_emit_plot_float(name, val);
}

void GoTracyPlotDoublet(char *name, double val)
{
    ___tracy_emit_plot(name, val);
}

void GoTracyPlotInt(char *name, int val)
{
    ___tracy_emit_plot_int(name, val);
}

void GoTracyMessageAppinfo(char *info)
{
    ___tracy_emit_message_appinfo(info, strlen(info));
}

