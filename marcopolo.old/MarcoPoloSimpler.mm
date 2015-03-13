<map version="0.9.0">
<!-- To view this file, download free mind mapping software FreeMind from http://freemind.sourceforge.net -->
<node CREATED="1377473937557" ID="ID_121007360" MODIFIED="1377480516299" TEXT="Marco Polo Simpler&#xa;Just simple marco / polo: &#xa; ie find another app," VGAP="37">
<node CREATED="1377479796196" ID="ID_402495146" MODIFIED="1377486256531" POSITION="right" TEXT="marco.polo.Xxx|JSONparam">
<node CREATED="1377473966721" FOLDED="true" ID="ID_857094091" MODIFIED="1377561393479" TEXT="marco.polo.regapp|">
<node CREATED="1377478044520" FOLDED="true" ID="ID_963332620" MODIFIED="1377481968613" TEXT="mutliple registers ??&#xa;ie if regapp with already regd name ?&#xa;(couldbe controlled w. appreg param)">
<icon BUILTIN="help"/>
<node CREATED="1377478254165" ID="ID_1162451181" MODIFIED="1377478300770" TEXT="fail  ">
<icon BUILTIN="help"/>
</node>
<node CREATED="1377478289910" FOLDED="true" ID="ID_1149773362" MODIFIED="1377481967631" TEXT="override ">
<icon BUILTIN="button_ok"/>
<node CREATED="1377481913635" ID="ID_1077575577" MODIFIED="1377481917867" TEXT="for now: override"/>
</node>
<node CREATED="1377478278142" ID="ID_1030696836" MODIFIED="1377478302020" TEXT="keep multi">
<icon BUILTIN="help"/>
</node>
</node>
<node CREATED="1377480160035" ID="ID_143595742" MODIFIED="1377480167974" TEXT="param appname"/>
</node>
<node CREATED="1377473988651" FOLDED="true" HGAP="23" ID="ID_1085942607" MODIFIED="1377561394819" TEXT="marco.polo.unregapp|" VSHIFT="17">
<node CREATED="1377480169740" ID="ID_1598541055" MODIFIED="1377480170661" TEXT="param appname"/>
</node>
<node CREATED="1377474013807" FOLDED="true" HGAP="23" ID="ID_1967408645" MODIFIED="1377561397730" TEXT="marco.polo.qryapp|" VSHIFT="16">
<node CREATED="1377480173847" ID="ID_286851933" MODIFIED="1377480174942" TEXT="param appname"/>
<node CREATED="1377474495936" FOLDED="true" ID="ID_378677923" MODIFIED="1377481852142" TEXT="marcoPolo will &apos;ping&apos; the app&#xa;">
<node CREATED="1377478325893" ID="ID_223445232" MODIFIED="1377478345164" TEXT="if gets answer then&#xa;   send app address back to requester&#xa;"/>
<node CREATED="1377478335810" ID="ID_839467869" MODIFIED="1377478353075" TEXT="else, no answer&#xa;  unreg app&#xa;  send back &apos;no app&apos; answer/err to requester"/>
</node>
</node>
</node>
<node CREATED="1377485944849" ID="ID_533838704" MODIFIED="1377486292118" POSITION="right" TEXT="note: prefer &apos;|&apos; sep versus &apos;:&apos;&#xa;to try to make safer (?) and differ from, for eg, a URL&#xa;like &quot;http://ab/cdef&quot;&#xa;"/>
<node CREATED="1377480229582" ID="ID_679687064" MODIFIED="1377481799799" POSITION="right" TEXT="JSON objects">
<node CREATED="1377479984981" ID="ID_492362965" MODIFIED="1377481839781" TEXT="stdHeader">
<node CREATED="1377479988848" FOLDED="true" ID="ID_154970330" MODIFIED="1377561106292" TEXT="version">
<node CREATED="1377479991522" ID="ID_1171275102" MODIFIED="1377479999569" TEXT="Major int"/>
<node CREATED="1377479993476" ID="ID_512176747" MODIFIED="1377480001840" TEXT="Minor int"/>
</node>
</node>
<node CREATED="1377480120235" ID="ID_589727849" MODIFIED="1377481826752" TEXT="params">
<node CREATED="1377480131464" ID="ID_453472200" MODIFIED="1377561107431" TEXT="appname param">
<node CREATED="1377479967540" ID="ID_1016459185" MODIFIED="1377479972617" TEXT="StdHeader"/>
<node CREATED="1377479973201" ID="ID_1045589089" MODIFIED="1377479979200" TEXT="appname string"/>
</node>
</node>
<node CREATED="1377480253905" ID="ID_1638749447" MODIFIED="1377481827702" TEXT="answers">
<node CREATED="1377480257779" ID="ID_1435220336" MODIFIED="1377561108732" TEXT="stdAnswer">
<node CREATED="1377480272699" ID="ID_146917791" MODIFIED="1377480275911" TEXT="stdHeader"/>
<node CREATED="1377480276420" ID="ID_338499475" MODIFIED="1377480281637" TEXT="ok bool"/>
<node CREATED="1377480297386" ID="ID_1705882395" MODIFIED="1377480301537" TEXT="value string"/>
<node CREATED="1377480284464" ID="ID_481024568" MODIFIED="1377480288334" TEXT="error string"/>
</node>
</node>
</node>
<node CREATED="1377101722881" FOLDED="true" HGAP="47" ID="ID_1984794553" MODIFIED="1377480422285" POSITION="left" TEXT="msg name convention" VSHIFT="26">
<node CREATED="1377101733263" HGAP="22" ID="ID_1176621873" MODIFIED="1377101892602" TEXT="Domain&#xa;eg kwez.org" VSHIFT="5">
<icon BUILTIN="full-1"/>
</node>
<node CREATED="1377101754511" ID="ID_124424908" MODIFIED="1377101901955" TEXT="AppName&#xa;eg androidPush">
<icon BUILTIN="full-2"/>
</node>
<node CREATED="1377101767994" ID="ID_763250280" MODIFIED="1377101913513" TEXT="Method/msg etc&#xa;eg notifPushAvail">
<icon BUILTIN="full-3"/>
</node>
<node CREATED="1377101797478" HGAP="27" ID="ID_1615771619" MODIFIED="1377132873116" TEXT="&apos;kwez.org/androidPush/notifPushAvail&apos;" VSHIFT="19">
<icon BUILTIN="info"/>
</node>
</node>
<node CREATED="1377474213468" FOLDED="true" HGAP="-6" ID="ID_312475021" MODIFIED="1377481821108" POSITION="left" TEXT="the address of the udp connection used by/when app registers&#xa;is kept as the address of the app.&#xa;the app must listen on it for :&#xa;" VSHIFT="27">
<node CREATED="1377478380522" HGAP="28" ID="ID_1650652808" MODIFIED="1377481819168" TEXT="marco.polo.ping -&gt; answer &apos;ok&apos;&#xa;" VSHIFT="10"/>
<node CREATED="1377478395349" ID="ID_75970215" MODIFIED="1377478395350" TEXT="msgs sent by other apps"/>
</node>
<node CREATED="1377478156312" ID="ID_1396776112" MODIFIED="1377478235818" POSITION="left" TEXT="possible broadcast msgs">
<icon BUILTIN="help"/>
</node>
<node CREATED="1377478203849" ID="ID_143901562" MODIFIED="1377478232413" POSITION="left" TEXT="possible periodic pings&#xa;to keep list clean">
<icon BUILTIN="help"/>
</node>
</node>
</map>
