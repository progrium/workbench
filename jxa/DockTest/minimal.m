#import <Cocoa/Cocoa.h>
#include <stdlib.h>
int main ()
{
    NSString *self = [NSProcessInfo.processInfo.arguments[0] stringByAppendingString:@"2"];
    const char *path = [self cStringUsingEncoding:NSASCIIStringEncoding];
    execl(path, path, NULL);
    return 0;
}