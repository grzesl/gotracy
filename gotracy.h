#ifndef __MAIN_H__
#define __MAIN_H__

#include <windows.h>
#include <stdint.h>
#include <string.h>

#ifdef __cplusplus
extern "C"
{
#endif
void GoTracySetThreadName(const char*name);

int GoTracyZoneBegin(const char*name,const char *function,const char*file, uint32_t line, uint32_t color);
void GoTracyZoneEnd(int c);
void GoTracyZoneValue(int c, uint64_t value);
void GoTracyZoneText(int c, char* text);

void GoTracyMessageL(char * msg);
void GoTracyMessageLC(char * msg, uint32_t color);

void GoTracyFrameMark();
void GoTracyFrameMarkName(char *name);
void GoTracyFrameMarkStart(char *name);
void GoTracyFrameMarkEnd(char *name);

void GoTracyPlotFloat(char *name, float val);
void GoTracyPlotDoublet(char *name, double val);
void GoTracyPlotInt(char *name, int val);
void GoTracyMessageAppinfo(char *info);

void GoTracyMemoryAlloc(unsigned long long ptr, size_t size, int secure );
void GoTracyMemoryAllocNamed( unsigned long long ptr, size_t size, int secure, const char* name);
void GoTracyMemoryFree(unsigned long long ptr, int secure );
void GoTracyMemoryFreeNamed(unsigned long longptr, int secure, const char* name );

#ifdef __cplusplus
}
#endif

#endif // __MAIN_H__
