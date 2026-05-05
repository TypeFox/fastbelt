# Runtime ATN for grammar

## Grammar

```mermaid
flowchart TD
    q0(["SN:0<br/>RuleStart"])
    q1(["SN:1<br/>RuleStop"])
    q248["SN:248<br/>Basic<br/>"]
    q249["SN:249<br/>Basic<br/>"]
    q250["SN:250<br/>Basic<br/>"]
    q251["SN:251<br/>Basic<br/>"]
    q252["SN:252<br/>Basic<br/>"]
    q253["SN:253<br/>Basic<br/>"]
    q254["SN:254<br/>Basic<br/>"]
    q255["SN:255<br/>Basic<br/>"]
    q256["SN:256<br/>Basic<br/>"]
    q257["SN:257<br/>Basic<br/>"]
    q258["SN:258<br/>Basic<br/>"]
    q259["SN:259<br/>Basic<br/>"]
    q260["SN:260<br/>Basic<br/>"]
    q261["SN:261<br/>Basic<br/>"]
    q262["SN:262<br/>Basic<br/>"]
    q263["SN:263<br/>BlockEnd<br/>"]
    q264{"SN:264<br/>StarLoopEntry<br/><br/>dec=5"}
    q265["SN:265<br/>LoopEnd<br/>"]
    q266["SN:266<br/>StarLoopBack<br/>"]

    q0 --> q248
    q248 -->|"tok(&quot;grammar&quot;)"| q249
    q249 --> q250
    q250 -->|"tok(ID)"| q251
    q251 --> q252
    q252 -->|"tok(&quot;;&quot;)"| q253
    q253 --> q264
    q254 -.->|"[ParserRule]"| q255
    q255 --> q263
    q256 -.->|"[Token]"| q257
    q257 --> q263
    q258 -.->|"[Interface]"| q259
    q259 --> q263
    q260 -.->|"[CompositeRule]"| q261
    q261 --> q263
    q262 --> q254
    q262 --> q256
    q262 --> q258
    q262 --> q260
    q263 --> q266
    q264 --> q262
    q264 --> q265
    q265 --> q1
    q266 --> q264
```

## Interface

```mermaid
flowchart TD
    q2(["SN:2<br/>RuleStart"])
    q3(["SN:3<br/>RuleStop"])
    q267["SN:267<br/>Basic<br/>"]
    q268["SN:268<br/>Basic<br/>"]
    q269["SN:269<br/>Basic<br/>"]
    q270["SN:270<br/>Basic<br/>"]
    q271["SN:271<br/>Basic<br/>"]
    q272["SN:272<br/>Basic<br/>"]
    q273["SN:273<br/>Basic<br/>"]
    q274["SN:274<br/>Basic<br/>"]
    q275["SN:275<br/>Basic<br/>"]
    q276["SN:276<br/>Basic<br/>"]
    q277["SN:277<br/>Basic<br/>"]
    q278["SN:278<br/>Basic<br/>"]
    q279{"SN:279<br/>StarLoopEntry<br/><br/>dec=6"}
    q280["SN:280<br/>LoopEnd<br/>"]
    q281["SN:281<br/>StarLoopBack<br/>"]
    q282["SN:282<br/>Basic<br/>"]
    q283["SN:283<br/>Basic<br/>"]
    q284["SN:284<br/>Basic<br/>"]
    q285["SN:285<br/>Basic<br/>"]
    q286{"SN:286<br/>StarLoopEntry<br/><br/>dec=7"}
    q287["SN:287<br/>LoopEnd<br/>"]
    q288["SN:288<br/>StarLoopBack<br/>"]
    q289["SN:289<br/>Basic<br/>"]
    q290["SN:290<br/>Basic<br/>"]

    q2 --> q267
    q267 -->|"tok(&quot;interface&quot;)"| q268
    q268 --> q269
    q269 -->|"tok(ID)"| q270
    q270 --> q271
    q271 -->|"tok(&quot;extends&quot;)"| q272
    q271 --> q280
    q272 --> q273
    q273 -->|"tok(ID)"| q274
    q274 --> q279
    q275 -->|"tok(&quot;,&quot;)"| q276
    q276 --> q277
    q277 -->|"tok(ID)"| q278
    q278 --> q281
    q279 --> q275
    q279 --> q280
    q280 --> q282
    q281 --> q279
    q282 -->|"tok(&quot;{&quot;)"| q283
    q283 --> q286
    q284 -.->|"[Field]"| q285
    q285 --> q288
    q286 --> q284
    q286 --> q287
    q287 --> q289
    q288 --> q286
    q289 -->|"tok(&quot;}&quot;)"| q290
    q290 --> q3
```

## Field

```mermaid
flowchart TD
    q4(["SN:4<br/>RuleStart"])
    q5(["SN:5<br/>RuleStop"])
    q116["SN:116<br/>Basic<br/>"]
    q117["SN:117<br/>Basic<br/>"]
    q118["SN:118<br/>Basic<br/>"]
    q119["SN:119<br/>Basic<br/>"]

    q4 --> q116
    q116 -->|"tok(ID)"| q117
    q117 --> q118
    q118 -.->|"[FieldType]"| q119
    q119 --> q5
```

## FieldType

```mermaid
flowchart TD
    q6(["SN:6<br/>RuleStart"])
    q7(["SN:7<br/>RuleStop"])
    q120["SN:120<br/>Basic<br/>"]
    q121["SN:121<br/>Basic<br/>"]
    q122["SN:122<br/>Basic<br/>"]
    q123["SN:123<br/>Basic<br/>"]
    q124["SN:124<br/>Basic<br/>"]
    q125["SN:125<br/>Basic<br/>"]
    q126["SN:126<br/>Basic<br/>"]
    q127["SN:127<br/>Basic<br/>"]
    q128["SN:128<br/>Basic<br/>"]
    q129["SN:129<br/>BlockEnd<br/>"]

    q6 --> q128
    q120 -.->|"[SimpleType]"| q121
    q121 --> q129
    q122 -.->|"[ReferenceType]"| q123
    q123 --> q129
    q124 -.->|"[ArrayType]"| q125
    q125 --> q129
    q126 -.->|"[PrimitiveType]"| q127
    q127 --> q129
    q128 --> q120
    q128 --> q122
    q128 --> q124
    q128 --> q126
    q129 --> q7
```

## ArrayType

```mermaid
flowchart TD
    q8(["SN:8<br/>RuleStart"])
    q9(["SN:9<br/>RuleStop"])
    q305["SN:305<br/>Basic<br/>"]
    q306["SN:306<br/>Basic<br/>"]
    q307["SN:307<br/>Basic<br/>"]
    q308["SN:308<br/>Basic<br/>"]
    q309["SN:309<br/>Basic<br/>"]
    q310["SN:310<br/>Basic<br/>"]

    q8 --> q305
    q305 -->|"tok(&quot;[&quot;)"| q306
    q306 --> q307
    q307 -->|"tok(&quot;]&quot;)"| q308
    q308 --> q309
    q309 -.->|"[FieldType]"| q310
    q310 --> q9
```

## ReferenceType

```mermaid
flowchart TD
    q10(["SN:10<br/>RuleStart"])
    q11(["SN:11<br/>RuleStop"])
    q174["SN:174<br/>Basic<br/>"]
    q175["SN:175<br/>Basic<br/>"]
    q176["SN:176<br/>Basic<br/>"]
    q177["SN:177<br/>Basic<br/>"]

    q10 --> q174
    q174 -->|"tok(&quot;*&quot;)"| q175
    q175 --> q176
    q176 -->|"tok(ID)"| q177
    q177 --> q11
```

## SimpleType

```mermaid
flowchart TD
    q12(["SN:12<br/>RuleStart"])
    q13(["SN:13<br/>RuleStop"])
    q291["SN:291<br/>Basic<br/>"]
    q292["SN:292<br/>Basic<br/>"]

    q12 --> q291
    q291 -->|"tok(ID)"| q292
    q292 --> q13
```

## PrimitiveType

```mermaid
flowchart TD
    q14(["SN:14<br/>RuleStart"])
    q15(["SN:15<br/>RuleStop"])
    q178["SN:178<br/>Basic<br/>"]
    q179["SN:179<br/>Basic<br/>"]
    q180["SN:180<br/>Basic<br/>"]
    q181["SN:181<br/>Basic<br/>"]
    q182["SN:182<br/>Basic<br/>"]
    q183["SN:183<br/>Basic<br/>"]
    q184["SN:184<br/>Basic<br/>"]
    q185["SN:185<br/>BlockEnd<br/>"]

    q14 --> q184
    q178 -->|"tok(&quot;string&quot;)"| q179
    q179 --> q185
    q180 -->|"tok(&quot;bool&quot;)"| q181
    q181 --> q185
    q182 -->|"tok(&quot;composite&quot;)"| q183
    q183 --> q185
    q184 --> q178
    q184 --> q180
    q184 --> q182
    q185 --> q15
```

## ParserRule

```mermaid
flowchart TD
    q16(["SN:16<br/>RuleStart"])
    q17(["SN:17<br/>RuleStop"])
    q293["SN:293<br/>Basic<br/>"]
    q294["SN:294<br/>Basic<br/>"]
    q295["SN:295<br/>Basic<br/>"]
    q296["SN:296<br/>Basic<br/>"]
    q297["SN:297<br/>Basic<br/>"]
    q298["SN:298<br/>Basic<br/>"]
    q299["SN:299<br/>Basic<br/>"]
    q300["SN:300<br/>Basic<br/>"]
    q301["SN:301<br/>Basic<br/>"]
    q302["SN:302<br/>Basic<br/>"]
    q303["SN:303<br/>Basic<br/>"]
    q304["SN:304<br/>Basic<br/>"]

    q16 --> q293
    q293 -->|"tok(ID)"| q294
    q294 --> q295
    q295 -->|"tok(&quot;returns&quot;)"| q296
    q295 --> q298
    q296 --> q297
    q297 -->|"tok(ID)"| q298
    q298 --> q299
    q299 -->|"tok(&quot;:&quot;)"| q300
    q300 --> q301
    q301 -.->|"[Alternatives]"| q302
    q302 --> q303
    q303 -->|"tok(&quot;;&quot;)"| q304
    q304 --> q17
```

## Token

```mermaid
flowchart TD
    q18(["SN:18<br/>RuleStart"])
    q19(["SN:19<br/>RuleStop"])
    q50["SN:50<br/>Basic<br/>"]
    q51["SN:51<br/>Basic<br/>"]
    q52["SN:52<br/>Basic<br/>"]
    q53["SN:53<br/>Basic<br/>"]
    q54["SN:54<br/>Basic<br/>"]
    q55["SN:55<br/>BlockEnd<br/>"]
    q56["SN:56<br/>Basic<br/>"]
    q57["SN:57<br/>Basic<br/>"]
    q58["SN:58<br/>Basic<br/>"]
    q59["SN:59<br/>Basic<br/>"]
    q60["SN:60<br/>Basic<br/>"]
    q61["SN:61<br/>Basic<br/>"]
    q62["SN:62<br/>Basic<br/>"]
    q63["SN:63<br/>Basic<br/>"]
    q64["SN:64<br/>Basic<br/>"]
    q65["SN:65<br/>Basic<br/>"]

    q18 --> q54
    q50 -->|"tok(&quot;hidden&quot;)"| q51
    q51 --> q55
    q52 -->|"tok(&quot;comment&quot;)"| q53
    q53 --> q55
    q54 --> q50
    q54 --> q52
    q54 --> q55
    q55 --> q56
    q56 -->|"tok(&quot;token&quot;)"| q57
    q57 --> q58
    q58 -->|"tok(ID)"| q59
    q59 --> q60
    q60 -->|"tok(&quot;:&quot;)"| q61
    q61 --> q62
    q62 -->|"tok(RegexLiteral)"| q63
    q63 --> q64
    q64 -->|"tok(&quot;;&quot;)"| q65
    q65 --> q19
```

## Alternatives

```mermaid
flowchart TD
    q20(["SN:20<br/>RuleStart"])
    q21(["SN:21<br/>RuleStop"])
    q66["SN:66<br/>Basic<br/>"]
    q67["SN:67<br/>Basic<br/>"]
    q68["SN:68<br/>Basic<br/>"]
    q69["SN:69<br/>Basic<br/>"]
    q70["SN:70<br/>Basic<br/>"]
    q71["SN:71<br/>Basic<br/>"]
    q72{"SN:72<br/>PlusLoopBack<br/><br/>dec=0"}
    q73["SN:73<br/>LoopEnd<br/>"]

    q20 --> q66
    q66 -.->|"[Group]"| q67
    q67 --> q68
    q68 -->|"tok(&quot;|&quot;)"| q69
    q68 --> q73
    q69 --> q70
    q70 -.->|"[Group]"| q71
    q71 --> q72
    q72 --> q68
    q72 --> q73
    q73 --> q21
```

## Group

```mermaid
flowchart TD
    q22(["SN:22<br/>RuleStart"])
    q23(["SN:23<br/>RuleStop"])
    q186["SN:186<br/>Basic<br/>"]
    q187["SN:187<br/>Basic<br/>"]
    q188["SN:188<br/>Basic<br/>"]
    q189["SN:189<br/>Basic<br/>"]
    q190{"SN:190<br/>PlusLoopBack<br/><br/>dec=3"}
    q191["SN:191<br/>LoopEnd<br/>"]

    q22 --> q186
    q186 -.->|"[Element]"| q187
    q187 --> q188
    q188 -.->|"[Element]"| q189
    q188 --> q191
    q189 --> q190
    q190 --> q188
    q190 --> q191
    q191 --> q23
```

## Element

```mermaid
flowchart TD
    q24(["SN:24<br/>RuleStart"])
    q25(["SN:25<br/>RuleStop"])
    q130["SN:130<br/>Basic<br/>"]
    q131["SN:131<br/>Basic<br/>"]
    q132["SN:132<br/>Basic<br/>"]
    q133["SN:133<br/>Basic<br/>"]
    q134["SN:134<br/>Basic<br/>"]
    q135["SN:135<br/>Basic<br/>"]
    q136["SN:136<br/>Basic<br/>"]
    q137["SN:137<br/>Basic<br/>"]
    q138["SN:138<br/>Basic<br/>"]
    q139["SN:139<br/>Basic<br/>"]
    q140["SN:140<br/>Basic<br/>"]
    q141["SN:141<br/>Basic<br/>"]
    q142["SN:142<br/>Basic<br/>"]
    q143["SN:143<br/>Basic<br/>"]
    q144["SN:144<br/>Basic<br/>"]
    q145["SN:145<br/>BlockEnd<br/>"]
    q146["SN:146<br/>Basic<br/>"]
    q147["SN:147<br/>Basic<br/>"]
    q148["SN:148<br/>Basic<br/>"]
    q149["SN:149<br/>Basic<br/>"]
    q150["SN:150<br/>Basic<br/>"]
    q151["SN:151<br/>Basic<br/>"]
    q152["SN:152<br/>Basic<br/>"]
    q153["SN:153<br/>BlockEnd<br/>"]

    q24 --> q144
    q130 -.->|"[Keyword]"| q131
    q131 --> q145
    q132 -.->|"[Assignment]"| q133
    q133 --> q145
    q134 -.->|"[RuleCall]"| q135
    q135 --> q145
    q136 -.->|"[Action]"| q137
    q137 --> q145
    q138 -->|"tok(&quot;(&quot;)"| q139
    q139 --> q140
    q140 -.->|"[Alternatives]"| q141
    q141 --> q142
    q142 -->|"tok(&quot;)&quot;)"| q143
    q143 --> q145
    q144 --> q130
    q144 --> q132
    q144 --> q134
    q144 --> q136
    q144 --> q138
    q145 --> q152
    q146 -->|"tok(&quot;*&quot;)"| q147
    q147 --> q153
    q148 -->|"tok(&quot;+&quot;)"| q149
    q149 --> q153
    q150 -->|"tok(&quot;?&quot;)"| q151
    q151 --> q153
    q152 --> q146
    q152 --> q148
    q152 --> q150
    q152 --> q153
    q153 --> q25
```

## Keyword

```mermaid
flowchart TD
    q26(["SN:26<br/>RuleStart"])
    q27(["SN:27<br/>RuleStop"])
    q154["SN:154<br/>Basic<br/>"]
    q155["SN:155<br/>Basic<br/>"]

    q26 --> q154
    q154 -->|"tok(StringLiteral)"| q155
    q155 --> q27
```

## Assignment

```mermaid
flowchart TD
    q28(["SN:28<br/>RuleStart"])
    q29(["SN:29<br/>RuleStop"])
    q74["SN:74<br/>Basic<br/>"]
    q75["SN:75<br/>Basic<br/>"]
    q76["SN:76<br/>Basic<br/>"]
    q77["SN:77<br/>Basic<br/>"]
    q78["SN:78<br/>Basic<br/>"]
    q79["SN:79<br/>Basic<br/>"]
    q80["SN:80<br/>Basic<br/>"]
    q81["SN:81<br/>Basic<br/>"]
    q82["SN:82<br/>Basic<br/>"]
    q83["SN:83<br/>BlockEnd<br/>"]
    q84["SN:84<br/>Basic<br/>"]
    q85["SN:85<br/>Basic<br/>"]

    q28 --> q74
    q74 -->|"tok(ID)"| q75
    q75 --> q82
    q76 -->|"tok(&quot;+=&quot;)"| q77
    q77 --> q83
    q78 -->|"tok(&quot;=&quot;)"| q79
    q79 --> q83
    q80 -->|"tok(&quot;?=&quot;)"| q81
    q81 --> q83
    q82 --> q76
    q82 --> q78
    q82 --> q80
    q83 --> q84
    q84 -.->|"[Assignable]"| q85
    q85 --> q29
```

## Assignable

```mermaid
flowchart TD
    q30(["SN:30<br/>RuleStart"])
    q31(["SN:31<br/>RuleStop"])
    q192["SN:192<br/>Basic<br/>"]
    q193["SN:193<br/>Basic<br/>"]
    q194["SN:194<br/>Basic<br/>"]
    q195["SN:195<br/>Basic<br/>"]
    q196["SN:196<br/>Basic<br/>"]
    q197["SN:197<br/>Basic<br/>"]
    q198["SN:198<br/>Basic<br/>"]
    q199["SN:199<br/>Basic<br/>"]
    q200["SN:200<br/>Basic<br/>"]
    q201["SN:201<br/>Basic<br/>"]
    q202["SN:202<br/>Basic<br/>"]
    q203["SN:203<br/>Basic<br/>"]
    q204["SN:204<br/>Basic<br/>"]
    q205["SN:205<br/>BlockEnd<br/>"]

    q30 --> q204
    q192 -.->|"[Keyword]"| q193
    q193 --> q205
    q194 -.->|"[RuleCall]"| q195
    q195 --> q205
    q196 -.->|"[CrossRef]"| q197
    q197 --> q205
    q198 -->|"tok(&quot;(&quot;)"| q199
    q199 --> q200
    q200 -.->|"[AssignableAlternatives]"| q201
    q201 --> q202
    q202 -->|"tok(&quot;)&quot;)"| q203
    q203 --> q205
    q204 --> q192
    q204 --> q194
    q204 --> q196
    q204 --> q198
    q205 --> q31
```

## AssignableWithoutAlts

```mermaid
flowchart TD
    q32(["SN:32<br/>RuleStart"])
    q33(["SN:33<br/>RuleStop"])
    q206["SN:206<br/>Basic<br/>"]
    q207["SN:207<br/>Basic<br/>"]
    q208["SN:208<br/>Basic<br/>"]
    q209["SN:209<br/>Basic<br/>"]
    q210["SN:210<br/>Basic<br/>"]
    q211["SN:211<br/>Basic<br/>"]
    q212["SN:212<br/>Basic<br/>"]
    q213["SN:213<br/>BlockEnd<br/>"]

    q32 --> q212
    q206 -.->|"[Keyword]"| q207
    q207 --> q213
    q208 -.->|"[RuleCall]"| q209
    q209 --> q213
    q210 -.->|"[CrossRef]"| q211
    q211 --> q213
    q212 --> q206
    q212 --> q208
    q212 --> q210
    q213 --> q33
```

## AssignableAlternatives

```mermaid
flowchart TD
    q34(["SN:34<br/>RuleStart"])
    q35(["SN:35<br/>RuleStop"])
    q156["SN:156<br/>Basic<br/>"]
    q157["SN:157<br/>Basic<br/>"]
    q158["SN:158<br/>Basic<br/>"]
    q159["SN:159<br/>Basic<br/>"]
    q160["SN:160<br/>Basic<br/>"]
    q161["SN:161<br/>Basic<br/>"]
    q162{"SN:162<br/>PlusLoopBack<br/><br/>dec=2"}
    q163["SN:163<br/>LoopEnd<br/>"]

    q34 --> q156
    q156 -.->|"[AssignableWithoutAlts]"| q157
    q157 --> q158
    q158 -->|"tok(&quot;|&quot;)"| q159
    q158 --> q163
    q159 --> q160
    q160 -.->|"[AssignableWithoutAlts]"| q161
    q161 --> q162
    q162 --> q158
    q162 --> q163
    q163 --> q35
```

## CrossRef

```mermaid
flowchart TD
    q36(["SN:36<br/>RuleStart"])
    q37(["SN:37<br/>RuleStop"])
    q214["SN:214<br/>Basic<br/>"]
    q215["SN:215<br/>Basic<br/>"]
    q216["SN:216<br/>Basic<br/>"]
    q217["SN:217<br/>Basic<br/>"]
    q218["SN:218<br/>Basic<br/>"]
    q219["SN:219<br/>Basic<br/>"]
    q220["SN:220<br/>Basic<br/>"]
    q221["SN:221<br/>Basic<br/>"]
    q222["SN:222<br/>Basic<br/>"]
    q223["SN:223<br/>Basic<br/>"]

    q36 --> q214
    q214 -->|"tok(&quot;[&quot;)"| q215
    q215 --> q216
    q216 -->|"tok(ID)"| q217
    q217 --> q218
    q218 -->|"tok(&quot;:&quot;)"| q219
    q218 --> q221
    q219 --> q220
    q220 -.->|"[RuleCall]"| q221
    q221 --> q222
    q222 -->|"tok(&quot;]&quot;)"| q223
    q223 --> q37
```

## RuleCall

```mermaid
flowchart TD
    q38(["SN:38<br/>RuleStart"])
    q39(["SN:39<br/>RuleStop"])
    q86["SN:86<br/>Basic<br/>"]
    q87["SN:87<br/>Basic<br/>"]

    q38 --> q86
    q86 -->|"tok(ID)"| q87
    q87 --> q39
```

## Action

```mermaid
flowchart TD
    q40(["SN:40<br/>RuleStart"])
    q41(["SN:41<br/>RuleStop"])
    q224["SN:224<br/>Basic<br/>"]
    q225["SN:225<br/>Basic<br/>"]
    q226["SN:226<br/>Basic<br/>"]
    q227["SN:227<br/>Basic<br/>"]
    q228["SN:228<br/>Basic<br/>"]
    q229["SN:229<br/>Basic<br/>"]
    q230["SN:230<br/>Basic<br/>"]
    q231["SN:231<br/>Basic<br/>"]
    q232["SN:232<br/>Basic<br/>"]
    q233["SN:233<br/>Basic<br/>"]
    q234["SN:234<br/>Basic<br/>"]
    q235["SN:235<br/>Basic<br/>"]
    q236["SN:236<br/>Basic<br/>"]
    q237["SN:237<br/>BlockEnd<br/>"]
    q238["SN:238<br/>Basic<br/>"]
    q239["SN:239<br/>Basic<br/>"]
    q240["SN:240<br/>Basic<br/>"]
    q241["SN:241<br/>Basic<br/>"]

    q40 --> q224
    q224 -->|"tok(&quot;{&quot;)"| q225
    q225 --> q226
    q226 -->|"tok(ID)"| q227
    q227 --> q228
    q228 -->|"tok(&quot;.&quot;)"| q229
    q228 --> q239
    q229 --> q230
    q230 -->|"tok(ID)"| q231
    q231 --> q236
    q232 -->|"tok(&quot;+=&quot;)"| q233
    q233 --> q237
    q234 -->|"tok(&quot;=&quot;)"| q235
    q235 --> q237
    q236 --> q232
    q236 --> q234
    q237 --> q238
    q238 -->|"tok(&quot;current&quot;)"| q239
    q239 --> q240
    q240 -->|"tok(&quot;}&quot;)"| q241
    q241 --> q41
```

## CompositeRule

```mermaid
flowchart TD
    q42(["SN:42<br/>RuleStart"])
    q43(["SN:43<br/>RuleStop"])
    q164["SN:164<br/>Basic<br/>"]
    q165["SN:165<br/>Basic<br/>"]
    q166["SN:166<br/>Basic<br/>"]
    q167["SN:167<br/>Basic<br/>"]
    q168["SN:168<br/>Basic<br/>"]
    q169["SN:169<br/>Basic<br/>"]
    q170["SN:170<br/>Basic<br/>"]
    q171["SN:171<br/>Basic<br/>"]
    q172["SN:172<br/>Basic<br/>"]
    q173["SN:173<br/>Basic<br/>"]

    q42 --> q164
    q164 -->|"tok(&quot;composite&quot;)"| q165
    q165 --> q166
    q166 -->|"tok(ID)"| q167
    q167 --> q168
    q168 -->|"tok(&quot;:&quot;)"| q169
    q169 --> q170
    q170 -.->|"[CompositeAlternatives]"| q171
    q171 --> q172
    q172 -->|"tok(&quot;;&quot;)"| q173
    q173 --> q43
```

## CompositeAlternatives

```mermaid
flowchart TD
    q44(["SN:44<br/>RuleStart"])
    q45(["SN:45<br/>RuleStop"])
    q88["SN:88<br/>Basic<br/>"]
    q89["SN:89<br/>Basic<br/>"]
    q90["SN:90<br/>Basic<br/>"]
    q91["SN:91<br/>Basic<br/>"]
    q92["SN:92<br/>Basic<br/>"]
    q93["SN:93<br/>Basic<br/>"]
    q94{"SN:94<br/>PlusLoopBack<br/><br/>dec=1"}
    q95["SN:95<br/>LoopEnd<br/>"]

    q44 --> q88
    q88 -.->|"[CompositeGroup]"| q89
    q89 --> q90
    q90 -->|"tok(&quot;|&quot;)"| q91
    q90 --> q95
    q91 --> q92
    q92 -.->|"[CompositeGroup]"| q93
    q93 --> q94
    q94 --> q90
    q94 --> q95
    q95 --> q45
```

## CompositeGroup

```mermaid
flowchart TD
    q46(["SN:46<br/>RuleStart"])
    q47(["SN:47<br/>RuleStop"])
    q242["SN:242<br/>Basic<br/>"]
    q243["SN:243<br/>Basic<br/>"]
    q244["SN:244<br/>Basic<br/>"]
    q245["SN:245<br/>Basic<br/>"]
    q246{"SN:246<br/>PlusLoopBack<br/><br/>dec=4"}
    q247["SN:247<br/>LoopEnd<br/>"]

    q46 --> q242
    q242 -.->|"[CompositeElement]"| q243
    q243 --> q244
    q244 -.->|"[CompositeElement]"| q245
    q244 --> q247
    q245 --> q246
    q246 --> q244
    q246 --> q247
    q247 --> q47
```

## CompositeElement

```mermaid
flowchart TD
    q48(["SN:48<br/>RuleStart"])
    q49(["SN:49<br/>RuleStop"])
    q96["SN:96<br/>Basic<br/>"]
    q97["SN:97<br/>Basic<br/>"]
    q98["SN:98<br/>Basic<br/>"]
    q99["SN:99<br/>Basic<br/>"]
    q100["SN:100<br/>Basic<br/>"]
    q101["SN:101<br/>Basic<br/>"]
    q102["SN:102<br/>Basic<br/>"]
    q103["SN:103<br/>Basic<br/>"]
    q104["SN:104<br/>Basic<br/>"]
    q105["SN:105<br/>Basic<br/>"]
    q106["SN:106<br/>Basic<br/>"]
    q107["SN:107<br/>BlockEnd<br/>"]
    q108["SN:108<br/>Basic<br/>"]
    q109["SN:109<br/>Basic<br/>"]
    q110["SN:110<br/>Basic<br/>"]
    q111["SN:111<br/>Basic<br/>"]
    q112["SN:112<br/>Basic<br/>"]
    q113["SN:113<br/>Basic<br/>"]
    q114["SN:114<br/>Basic<br/>"]
    q115["SN:115<br/>BlockEnd<br/>"]

    q48 --> q106
    q96 -.->|"[Keyword]"| q97
    q97 --> q107
    q98 -.->|"[RuleCall]"| q99
    q99 --> q107
    q100 -->|"tok(&quot;(&quot;)"| q101
    q101 --> q102
    q102 -.->|"[CompositeAlternatives]"| q103
    q103 --> q104
    q104 -->|"tok(&quot;)&quot;)"| q105
    q105 --> q107
    q106 --> q96
    q106 --> q98
    q106 --> q100
    q107 --> q114
    q108 -->|"tok(&quot;*&quot;)"| q109
    q109 --> q115
    q110 -->|"tok(&quot;+&quot;)"| q111
    q111 --> q115
    q112 -->|"tok(&quot;?&quot;)"| q113
    q113 --> q115
    q114 --> q108
    q114 --> q110
    q114 --> q112
    q114 --> q115
    q115 --> q49
```

