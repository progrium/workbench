#import <Cocoa/Cocoa.h>

int main(int argc, char *argv[])
{
    return 0;
}

@interface DockTileDemo : NSObject <NSDockTilePlugIn>
@end

@implementation DockTileDemo
- (void)setDockTile:(NSDockTile *)dockTile {
    NSLog(@"setDockTile3");
    if (dockTile) {
        [dockTile setBadgeLabel:[NSString stringWithFormat:@"%i", 3]];
    }
    
}
@end