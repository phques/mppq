// Java Package provider is a proxy for talking to a Go program.
//   gobind -lang=java github.com/phques/mppq/examples/libprovider/provider
//
// File is generated by gobind. Do not edit.
package go.provider;

import go.Seq;

public abstract class Provider {
    private Provider() {} // uninstantiable
    
    public static void InitAppFilesDir(String appFilesDir_) throws Exception {
        go.Seq _in = new go.Seq();
        go.Seq _out = new go.Seq();
        _in.writeUTF16(appFilesDir_);
        Seq.send(DESCRIPTOR, CALL_InitAppFilesDir, _in, _out);
        String _err = _out.readUTF16();
        if (_err != null) {
            throw new Exception(_err);
        }
    }
    
    public static void Start() throws Exception {
        go.Seq _in = new go.Seq();
        go.Seq _out = new go.Seq();
        Seq.send(DESCRIPTOR, CALL_Start, _in, _out);
        String _err = _out.readUTF16();
        if (_err != null) {
            throw new Exception(_err);
        }
    }
    
    public static void StartHTTP() throws Exception {
        go.Seq _in = new go.Seq();
        go.Seq _out = new go.Seq();
        Seq.send(DESCRIPTOR, CALL_StartHTTP, _in, _out);
        String _err = _out.readUTF16();
        if (_err != null) {
            throw new Exception(_err);
        }
    }
    
    private static final int CALL_InitAppFilesDir = 1;
    private static final int CALL_Start = 2;
    private static final int CALL_StartHTTP = 3;
    private static final String DESCRIPTOR = "provider";
}