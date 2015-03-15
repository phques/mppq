// Java Package provider is a proxy for talking to a Go program.
//   gobind -lang=java github.com/phques/mppq/examples/libprovider/provider
//
// File is generated by gobind. Do not edit.
package go.provider;

import go.Seq;

public abstract class Provider {
    private Provider() {} // uninstantiable
    
    public static void InitAppFilesDir(String appFilesDir_) {
        go.Seq _in = new go.Seq();
        go.Seq _out = new go.Seq();
        _in.writeUTF16(appFilesDir_);
        Seq.send(DESCRIPTOR, CALL_InitAppFilesDir, _in, _out);
    }
    
    public static void Register(String serviceName) {
        go.Seq _in = new go.Seq();
        go.Seq _out = new go.Seq();
        _in.writeUTF16(serviceName);
        Seq.send(DESCRIPTOR, CALL_Register, _in, _out);
    }
    
    public static void Start() {
        go.Seq _in = new go.Seq();
        go.Seq _out = new go.Seq();
        Seq.send(DESCRIPTOR, CALL_Start, _in, _out);
    }
    
    private static final int CALL_InitAppFilesDir = 1;
    private static final int CALL_Register = 2;
    private static final int CALL_Start = 3;
    private static final String DESCRIPTOR = "provider";
}
