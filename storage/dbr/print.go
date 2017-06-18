// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dbr

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/errors"
)

// QueryBuilder assembles a query and returns the raw SQL without parameter
// substitution and the arguments.
type QueryBuilder interface {
	ToSQL() (string, Arguments, error)
}

type queryBuilder interface {
	toSQL(queryWriter) error
	// appendArgs appends the arguments to Arguments and returns them. If
	// argument `Arguments` is nil, allocates new bytes
	appendArgs(Arguments) (Arguments, error)
	hasBuildCache() bool
	writeBuildCache(sql []byte)
	// readBuildCache returns the cached SQL string including its place holders.
	readBuildCache() (sql []byte, args Arguments, err error)
}

// queryWriter at used to generate a query.
type queryWriter interface {
	WriteString(s string) (n int, err error)
	WriteRune(r rune) (n int, err error)
	WriteByte(c byte) error
	Write(p []byte) (n int, err error)
}

var _ queryWriter = (*backHole)(nil)

type backHole struct{} // TODO(CyS) just a temporary implementation. should get removed later

func (backHole) WriteString(s string) (n int, err error) { return }
func (backHole) WriteRune(r rune) (n int, err error)     { return }
func (backHole) WriteByte(c byte) error                  { return nil }
func (backHole) Write(p []byte) (n int, err error)       { return }

// toSQL generates the SQL string and its place holders. Takes care of caching
// and interpolation. It returns the string with placeholders and a slice of
// query arguments. With switched on interpolation, it only returns a string
// including the stringyfied arguments. With an enabled cache, the arguments
// gets regenerated each time a call to ToSQL happens.
func toSQL(b queryBuilder, isInterpolate bool) (string, Arguments, error) {
	var ipBuf *bytes.Buffer // ip = interpolate buffer
	if isInterpolate {
		ipBuf = bufferpool.Get()
		defer bufferpool.Put(ipBuf)
	}

	useCache := b.hasBuildCache()
	if useCache {
		sql, args, err := b.readBuildCache()
		if err != nil {
			return "", nil, errors.Wrap(err, "[dbr] toSQL.readBuildCache")
		}
		if sql != nil {
			if isInterpolate {
				err := interpolate(ipBuf, sql, args...)
				return ipBuf.String(), nil, errors.Wrap(err, "[dbr] toSQL.Interpolate")
			}
			return string(sql), args, nil
		}
	}

	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	if err := b.toSQL(buf); err != nil {
		return "", nil, errors.Wrap(err, "[dbr] toSQL.toSQL")
	}
	// capacity of Arguments gets handled in the concret implementation of `b`
	args, err := b.appendArgs(Arguments{})
	if err != nil {
		return "", nil, errors.Wrap(err, "[dbr] toSQL.appendArgs")
	}
	if useCache {
		sqlCopy := make([]byte, buf.Len())
		copy(sqlCopy, buf.Bytes())
		b.writeBuildCache(sqlCopy)
	}

	if isInterpolate {
		err := interpolate(ipBuf, buf.Bytes(), args...)
		return ipBuf.String(), nil, errors.Wrap(err, "[dbr] toSQL.Interpolate")
	}
	return buf.String(), args, nil
}

func toSQLPrepared(b queryBuilder) (string, error) {
	// TODO(CyS) implement build cache like the toSQL function. see above.
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	err := b.toSQL(buf)
	return buf.String(), errors.Wrap(err, "[dbr] toSQLPrepared.toSQL")
}

func makeSQL(b queryBuilder, isInterpolate bool) string {
	sRaw, _, err := toSQL(b, isInterpolate)
	if err != nil {
		return fmt.Sprintf("[dbr] ToSQL Error: %+v", err)
	}
	return sRaw
}

// String returns a string representing a preprocessed, interpolated, query.
// On error, the error gets printed. Fulfills interface fmt.Stringer.
func (b *Delete) String() string {
	return makeSQL(b, b.IsInterpolate)
}

// String returns a string representing a preprocessed, interpolated, query.
// On error, the error gets printed. Fulfills interface fmt.Stringer.
func (b *Insert) String() string {
	return makeSQL(b, b.IsInterpolate)
}

// String returns a string representing a preprocessed, interpolated, query.
// On error, the error gets printed. Fulfills interface fmt.Stringer.
func (b *Select) String() string {
	return makeSQL(b, b.IsInterpolate)
}

// String returns a string representing a preprocessed, interpolated, query.
// On error, the error gets printed. Fulfills interface fmt.Stringer.
func (b *Update) String() string {
	return makeSQL(b, b.IsInterpolate)
}

func sqlWriteUnionAll(w queryWriter, isAll bool, isIntersect bool, isExcept bool) {
	w.WriteByte('\n')
	switch {
	case isIntersect:
		w.WriteString("INTERSECT") // MariaDB >= 10.3
	case isExcept:
		w.WriteString("EXCEPT") // MariaDB >= 10.3
	default:
		w.WriteString("UNION")
		if isAll {
			w.WriteString(" ALL")
		}
	}
	w.WriteByte('\n')
}

func sqlWriteOrderBy(w queryWriter, orderBys aliases, br bool) {
	if len(orderBys) == 0 {
		return
	}
	brS := ' '
	if br {
		brS = '\n'
	}
	w.WriteRune(brS)
	w.WriteString("ORDER BY ")
	for i, c := range orderBys {
		if i > 0 {
			w.WriteString(", ")
		}
		c.FquoteAs(w)
		// TODO append arguments
	}
}

func sqlWriteLimitOffset(w queryWriter, limitValid bool, limitCount uint64, offsetValid bool, offsetCount uint64) {
	if limitValid {
		w.WriteString(" LIMIT ")
		writeUint64(w, limitCount)
	}
	if offsetValid {
		w.WriteString(" OFFSET ")
		writeUint64(w, offsetCount)
	}
}

func writeFloat64(w queryWriter, f float64) error {
	_, err := w.WriteString(strconv.FormatFloat(f, 'f', -1, 64))
	return err
}

func writeInt64(w queryWriter, i int64) (err error) {
	if i >= 0 && i < nSmallItoas {
		_, err = w.WriteString(smallItoa(i))
	} else {
		_, err = w.WriteString(strconv.FormatInt(i, 10))
	}
	return err
}

func writeUint64(w queryWriter, i uint64) (err error) {
	if i < nSmallItoas {
		_, err = w.WriteString(smallItoa(int64(i)))
	} else {
		_, err = w.WriteString(strconv.FormatUint(i, 10))
	}
	return err
}

// smallItoa returns the string for an i with 0 <= i < nSmallItoas.
func smallItoa(i int64) string {
	var off int64
	if i < 10 {
		off = 3
	}
	if i >= 10 && i < 100 {
		off = 2
	}
	if i >= 100 && i < 1000 {
		off = 1
	}
	return smallsItoaString[i*4+off : i*4+4]
}

const nSmallItoas = 2000

const smallsItoaString = "00000001000200030004000500060007000800090010001100120013001400150016001700180019" +
	"00200021002200230024002500260027002800290030003100320033003400350036003700380039" +
	"00400041004200430044004500460047004800490050005100520053005400550056005700580059" +
	"00600061006200630064006500660067006800690070007100720073007400750076007700780079" +
	"00800081008200830084008500860087008800890090009100920093009400950096009700980099" +
	"01000101010201030104010501060107010801090110011101120113011401150116011701180119" +
	"01200121012201230124012501260127012801290130013101320133013401350136013701380139" +
	"01400141014201430144014501460147014801490150015101520153015401550156015701580159" +
	"01600161016201630164016501660167016801690170017101720173017401750176017701780179" +
	"01800181018201830184018501860187018801890190019101920193019401950196019701980199" +
	"02000201020202030204020502060207020802090210021102120213021402150216021702180219" +
	"02200221022202230224022502260227022802290230023102320233023402350236023702380239" +
	"02400241024202430244024502460247024802490250025102520253025402550256025702580259" +
	"02600261026202630264026502660267026802690270027102720273027402750276027702780279" +
	"02800281028202830284028502860287028802890290029102920293029402950296029702980299" +
	"03000301030203030304030503060307030803090310031103120313031403150316031703180319" +
	"03200321032203230324032503260327032803290330033103320333033403350336033703380339" +
	"03400341034203430344034503460347034803490350035103520353035403550356035703580359" +
	"03600361036203630364036503660367036803690370037103720373037403750376037703780379" +
	"03800381038203830384038503860387038803890390039103920393039403950396039703980399" +
	"04000401040204030404040504060407040804090410041104120413041404150416041704180419" +
	"04200421042204230424042504260427042804290430043104320433043404350436043704380439" +
	"04400441044204430444044504460447044804490450045104520453045404550456045704580459" +
	"04600461046204630464046504660467046804690470047104720473047404750476047704780479" +
	"04800481048204830484048504860487048804890490049104920493049404950496049704980499" +
	"05000501050205030504050505060507050805090510051105120513051405150516051705180519" +
	"05200521052205230524052505260527052805290530053105320533053405350536053705380539" +
	"05400541054205430544054505460547054805490550055105520553055405550556055705580559" +
	"05600561056205630564056505660567056805690570057105720573057405750576057705780579" +
	"05800581058205830584058505860587058805890590059105920593059405950596059705980599" +
	"06000601060206030604060506060607060806090610061106120613061406150616061706180619" +
	"06200621062206230624062506260627062806290630063106320633063406350636063706380639" +
	"06400641064206430644064506460647064806490650065106520653065406550656065706580659" +
	"06600661066206630664066506660667066806690670067106720673067406750676067706780679" +
	"06800681068206830684068506860687068806890690069106920693069406950696069706980699" +
	"07000701070207030704070507060707070807090710071107120713071407150716071707180719" +
	"07200721072207230724072507260727072807290730073107320733073407350736073707380739" +
	"07400741074207430744074507460747074807490750075107520753075407550756075707580759" +
	"07600761076207630764076507660767076807690770077107720773077407750776077707780779" +
	"07800781078207830784078507860787078807890790079107920793079407950796079707980799" +
	"08000801080208030804080508060807080808090810081108120813081408150816081708180819" +
	"08200821082208230824082508260827082808290830083108320833083408350836083708380839" +
	"08400841084208430844084508460847084808490850085108520853085408550856085708580859" +
	"08600861086208630864086508660867086808690870087108720873087408750876087708780879" +
	"08800881088208830884088508860887088808890890089108920893089408950896089708980899" +
	"09000901090209030904090509060907090809090910091109120913091409150916091709180919" +
	"09200921092209230924092509260927092809290930093109320933093409350936093709380939" +
	"09400941094209430944094509460947094809490950095109520953095409550956095709580959" +
	"09600961096209630964096509660967096809690970097109720973097409750976097709780979" +
	"09800981098209830984098509860987098809890990099109920993099409950996099709980999" +
	"10001001100210031004100510061007100810091010101110121013101410151016101710181019" +
	"10201021102210231024102510261027102810291030103110321033103410351036103710381039" +
	"10401041104210431044104510461047104810491050105110521053105410551056105710581059" +
	"10601061106210631064106510661067106810691070107110721073107410751076107710781079" +
	"10801081108210831084108510861087108810891090109110921093109410951096109710981099" +
	"11001101110211031104110511061107110811091110111111121113111411151116111711181119" +
	"11201121112211231124112511261127112811291130113111321133113411351136113711381139" +
	"11401141114211431144114511461147114811491150115111521153115411551156115711581159" +
	"11601161116211631164116511661167116811691170117111721173117411751176117711781179" +
	"11801181118211831184118511861187118811891190119111921193119411951196119711981199" +
	"12001201120212031204120512061207120812091210121112121213121412151216121712181219" +
	"12201221122212231224122512261227122812291230123112321233123412351236123712381239" +
	"12401241124212431244124512461247124812491250125112521253125412551256125712581259" +
	"12601261126212631264126512661267126812691270127112721273127412751276127712781279" +
	"12801281128212831284128512861287128812891290129112921293129412951296129712981299" +
	"13001301130213031304130513061307130813091310131113121313131413151316131713181319" +
	"13201321132213231324132513261327132813291330133113321333133413351336133713381339" +
	"13401341134213431344134513461347134813491350135113521353135413551356135713581359" +
	"13601361136213631364136513661367136813691370137113721373137413751376137713781379" +
	"13801381138213831384138513861387138813891390139113921393139413951396139713981399" +
	"14001401140214031404140514061407140814091410141114121413141414151416141714181419" +
	"14201421142214231424142514261427142814291430143114321433143414351436143714381439" +
	"14401441144214431444144514461447144814491450145114521453145414551456145714581459" +
	"14601461146214631464146514661467146814691470147114721473147414751476147714781479" +
	"14801481148214831484148514861487148814891490149114921493149414951496149714981499" +
	"15001501150215031504150515061507150815091510151115121513151415151516151715181519" +
	"15201521152215231524152515261527152815291530153115321533153415351536153715381539" +
	"15401541154215431544154515461547154815491550155115521553155415551556155715581559" +
	"15601561156215631564156515661567156815691570157115721573157415751576157715781579" +
	"15801581158215831584158515861587158815891590159115921593159415951596159715981599" +
	"16001601160216031604160516061607160816091610161116121613161416151616161716181619" +
	"16201621162216231624162516261627162816291630163116321633163416351636163716381639" +
	"16401641164216431644164516461647164816491650165116521653165416551656165716581659" +
	"16601661166216631664166516661667166816691670167116721673167416751676167716781679" +
	"16801681168216831684168516861687168816891690169116921693169416951696169716981699" +
	"17001701170217031704170517061707170817091710171117121713171417151716171717181719" +
	"17201721172217231724172517261727172817291730173117321733173417351736173717381739" +
	"17401741174217431744174517461747174817491750175117521753175417551756175717581759" +
	"17601761176217631764176517661767176817691770177117721773177417751776177717781779" +
	"17801781178217831784178517861787178817891790179117921793179417951796179717981799" +
	"18001801180218031804180518061807180818091810181118121813181418151816181718181819" +
	"18201821182218231824182518261827182818291830183118321833183418351836183718381839" +
	"18401841184218431844184518461847184818491850185118521853185418551856185718581859" +
	"18601861186218631864186518661867186818691870187118721873187418751876187718781879" +
	"18801881188218831884188518861887188818891890189118921893189418951896189718981899" +
	"19001901190219031904190519061907190819091910191119121913191419151916191719181919" +
	"19201921192219231924192519261927192819291930193119321933193419351936193719381939" +
	"19401941194219431944194519461947194819491950195119521953195419551956195719581959" +
	"19601961196219631964196519661967196819691970197119721973197419751976197719781979" +
	"19801981198219831984198519861987198819891990199119921993199419951996199719981999"
