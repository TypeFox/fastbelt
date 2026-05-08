# Runtime ATN for grammar

## Grammar

```mermaid
flowchart TD
    q0(["SN:0<br/>RuleStart"])
    q1(["SN:1<br/>RuleStop"])
    q50["SN:50<br/>Basic<br/>"]
    q51["SN:51<br/>Basic<br/>"]
    q52["SN:52<br/>Basic<br/>"]
    q53["SN:53<br/>Basic<br/>"]
    q54["SN:54<br/>Basic<br/>"]
    q55["SN:55<br/>Basic<br/>"]
    q56["SN:56<br/>Basic<br/>"]
    q57["SN:57<br/>Basic<br/>"]
    q58["SN:58<br/>Basic<br/>"]
    q59["SN:59<br/>Basic<br/>"]
    q60["SN:60<br/>Basic<br/>"]
    q61["SN:61<br/>Basic<br/>"]
    q62["SN:62<br/>Basic<br/>"]
    q63["SN:63<br/>Basic<br/>"]
    q64["SN:64<br/>Basic<br/>"]
    q65["SN:65<br/>BlockEnd<br/>"]
    q66{"SN:66<br/>StarLoopEntry<br/><br/>dec=0"}
    q67["SN:67<br/>LoopEnd<br/>"]
    q68["SN:68<br/>StarLoopBack<br/>"]

    q0 --> q50
    q50 -->|"tok(&quot;grammar&quot;)"| q51
    q51 --> q52
    q52 -->|"tok(ID)"| q53
    q53 --> q54
    q54 -->|"tok(&quot;;&quot;)"| q55
    q55 --> q66
    q56 -.->|"[ParserRule]"| q57
    q57 --> q65
    q58 -.->|"[Token]"| q59
    q59 --> q65
    q60 -.->|"[Interface]"| q61
    q61 --> q65
    q62 -.->|"[CompositeRule]"| q63
    q63 --> q65
    q64 --> q56
    q64 --> q58
    q64 --> q60
    q64 --> q62
    q65 --> q68
    q66 --> q64
    q66 --> q67
    q67 --> q1
    q68 --> q66
```

## Interface

```mermaid
flowchart TD
    q2(["SN:2<br/>RuleStart"])
    q3(["SN:3<br/>RuleStop"])
    q69["SN:69<br/>Basic<br/>"]
    q70["SN:70<br/>Basic<br/>"]
    q71["SN:71<br/>Basic<br/>"]
    q72["SN:72<br/>Basic<br/>"]
    q73["SN:73<br/>Basic<br/>"]
    q74["SN:74<br/>Basic<br/>"]
    q75["SN:75<br/>Basic<br/>"]
    q76["SN:76<br/>Basic<br/>"]
    q77["SN:77<br/>Basic<br/>"]
    q78["SN:78<br/>Basic<br/>"]
    q79["SN:79<br/>Basic<br/>"]
    q80["SN:80<br/>Basic<br/>"]
    q81{"SN:81<br/>StarLoopEntry<br/><br/>dec=1"}
    q82["SN:82<br/>LoopEnd<br/>"]
    q83["SN:83<br/>StarLoopBack<br/>"]
    q84["SN:84<br/>Basic<br/>"]
    q85["SN:85<br/>Basic<br/>"]
    q86["SN:86<br/>Basic<br/>"]
    q87["SN:87<br/>Basic<br/>"]
    q88{"SN:88<br/>StarLoopEntry<br/><br/>dec=2"}
    q89["SN:89<br/>LoopEnd<br/>"]
    q90["SN:90<br/>StarLoopBack<br/>"]
    q91["SN:91<br/>Basic<br/>"]
    q92["SN:92<br/>Basic<br/>"]

    q2 --> q69
    q69 -->|"tok(&quot;interface&quot;)"| q70
    q70 --> q71
    q71 -->|"tok(ID)"| q72
    q72 --> q73
    q73 -->|"tok(&quot;extends&quot;)"| q74
    q73 --> q82
    q74 --> q75
    q75 -->|"tok(ID)"| q76
    q76 --> q81
    q77 -->|"tok(&quot;,&quot;)"| q78
    q78 --> q79
    q79 -->|"tok(ID)"| q80
    q80 --> q83
    q81 --> q77
    q81 --> q82
    q82 --> q84
    q83 --> q81
    q84 -->|"tok(&quot;{&quot;)"| q85
    q85 --> q88
    q86 -.->|"[Field]"| q87
    q87 --> q90
    q88 --> q86
    q88 --> q89
    q89 --> q91
    q90 --> q88
    q91 -->|"tok(&quot;}&quot;)"| q92
    q92 --> q3
```

## Field

```mermaid
flowchart TD
    q4(["SN:4<br/>RuleStart"])
    q5(["SN:5<br/>RuleStop"])
    q93["SN:93<br/>Basic<br/>"]
    q94["SN:94<br/>Basic<br/>"]
    q95["SN:95<br/>Basic<br/>"]
    q96["SN:96<br/>Basic<br/>"]

    q4 --> q93
    q93 -->|"tok(ID)"| q94
    q94 --> q95
    q95 -.->|"[FieldType]"| q96
    q96 --> q5
```

## FieldType

```mermaid
flowchart TD
    q6(["SN:6<br/>RuleStart"])
    q7(["SN:7<br/>RuleStop"])
    q97["SN:97<br/>Basic<br/>"]
    q98["SN:98<br/>Basic<br/>"]
    q99["SN:99<br/>Basic<br/>"]
    q100["SN:100<br/>Basic<br/>"]
    q101["SN:101<br/>Basic<br/>"]
    q102["SN:102<br/>Basic<br/>"]
    q103["SN:103<br/>Basic<br/>"]
    q104["SN:104<br/>Basic<br/>"]
    q105["SN:105<br/>Basic<br/>"]
    q106["SN:106<br/>BlockEnd<br/>"]

    q6 --> q105
    q97 -.->|"[SimpleType]"| q98
    q98 --> q106
    q99 -.->|"[ReferenceType]"| q100
    q100 --> q106
    q101 -.->|"[ArrayType]"| q102
    q102 --> q106
    q103 -.->|"[PrimitiveType]"| q104
    q104 --> q106
    q105 --> q97
    q105 --> q99
    q105 --> q101
    q105 --> q103
    q106 --> q7
```

## ArrayType

```mermaid
flowchart TD
    q8(["SN:8<br/>RuleStart"])
    q9(["SN:9<br/>RuleStop"])
    q107["SN:107<br/>Basic<br/>"]
    q108["SN:108<br/>Basic<br/>"]
    q109["SN:109<br/>Basic<br/>"]
    q110["SN:110<br/>Basic<br/>"]
    q111["SN:111<br/>Basic<br/>"]
    q112["SN:112<br/>Basic<br/>"]

    q8 --> q107
    q107 -->|"tok(&quot;[&quot;)"| q108
    q108 --> q109
    q109 -->|"tok(&quot;]&quot;)"| q110
    q110 --> q111
    q111 -.->|"[FieldType]"| q112
    q112 --> q9
```

## ReferenceType

```mermaid
flowchart TD
    q10(["SN:10<br/>RuleStart"])
    q11(["SN:11<br/>RuleStop"])
    q113["SN:113<br/>Basic<br/>"]
    q114["SN:114<br/>Basic<br/>"]
    q115["SN:115<br/>Basic<br/>"]
    q116["SN:116<br/>Basic<br/>"]

    q10 --> q113
    q113 -->|"tok(&quot;*&quot;)"| q114
    q114 --> q115
    q115 -->|"tok(ID)"| q116
    q116 --> q11
```

## SimpleType

```mermaid
flowchart TD
    q12(["SN:12<br/>RuleStart"])
    q13(["SN:13<br/>RuleStop"])
    q117["SN:117<br/>Basic<br/>"]
    q118["SN:118<br/>Basic<br/>"]

    q12 --> q117
    q117 -->|"tok(ID)"| q118
    q118 --> q13
```

## PrimitiveType

```mermaid
flowchart TD
    q14(["SN:14<br/>RuleStart"])
    q15(["SN:15<br/>RuleStop"])
    q119["SN:119<br/>Basic<br/>"]
    q120["SN:120<br/>Basic<br/>"]
    q121["SN:121<br/>Basic<br/>"]
    q122["SN:122<br/>Basic<br/>"]
    q123["SN:123<br/>Basic<br/>"]
    q124["SN:124<br/>Basic<br/>"]
    q125["SN:125<br/>Basic<br/>"]
    q126["SN:126<br/>BlockEnd<br/>"]

    q14 --> q125
    q119 -->|"tok(&quot;string&quot;)"| q120
    q120 --> q126
    q121 -->|"tok(&quot;bool&quot;)"| q122
    q122 --> q126
    q123 -->|"tok(&quot;composite&quot;)"| q124
    q124 --> q126
    q125 --> q119
    q125 --> q121
    q125 --> q123
    q126 --> q15
```

## ParserRule

```mermaid
flowchart TD
    q16(["SN:16<br/>RuleStart"])
    q17(["SN:17<br/>RuleStop"])
    q127["SN:127<br/>Basic<br/>"]
    q128["SN:128<br/>Basic<br/>"]
    q129["SN:129<br/>Basic<br/>"]
    q130["SN:130<br/>Basic<br/>"]
    q131["SN:131<br/>Basic<br/>"]
    q132["SN:132<br/>Basic<br/>"]
    q133["SN:133<br/>Basic<br/>"]
    q134["SN:134<br/>Basic<br/>"]
    q135["SN:135<br/>Basic<br/>"]
    q136["SN:136<br/>Basic<br/>"]
    q137["SN:137<br/>Basic<br/>"]
    q138["SN:138<br/>Basic<br/>"]

    q16 --> q127
    q127 -->|"tok(ID)"| q128
    q128 --> q129
    q129 -->|"tok(&quot;returns&quot;)"| q130
    q129 --> q132
    q130 --> q131
    q131 -->|"tok(ID)"| q132
    q132 --> q133
    q133 -->|"tok(&quot;:&quot;)"| q134
    q134 --> q135
    q135 -.->|"[Alternatives]"| q136
    q136 --> q137
    q137 -->|"tok(&quot;;&quot;)"| q138
    q138 --> q17
```

## Token

```mermaid
flowchart TD
    q18(["SN:18<br/>RuleStart"])
    q19(["SN:19<br/>RuleStop"])
    q139["SN:139<br/>Basic<br/>"]
    q140["SN:140<br/>Basic<br/>"]
    q141["SN:141<br/>Basic<br/>"]
    q142["SN:142<br/>Basic<br/>"]
    q143["SN:143<br/>Basic<br/>"]
    q144["SN:144<br/>BlockEnd<br/>"]
    q145["SN:145<br/>Basic<br/>"]
    q146["SN:146<br/>Basic<br/>"]
    q147["SN:147<br/>Basic<br/>"]
    q148["SN:148<br/>Basic<br/>"]
    q149["SN:149<br/>Basic<br/>"]
    q150["SN:150<br/>Basic<br/>"]
    q151["SN:151<br/>Basic<br/>"]
    q152["SN:152<br/>Basic<br/>"]
    q153["SN:153<br/>Basic<br/>"]
    q154["SN:154<br/>Basic<br/>"]

    q18 --> q143
    q139 -->|"tok(&quot;hidden&quot;)"| q140
    q140 --> q144
    q141 -->|"tok(&quot;comment&quot;)"| q142
    q142 --> q144
    q143 --> q139
    q143 --> q141
    q143 --> q144
    q144 --> q145
    q145 -->|"tok(&quot;token&quot;)"| q146
    q146 --> q147
    q147 -->|"tok(ID)"| q148
    q148 --> q149
    q149 -->|"tok(&quot;:&quot;)"| q150
    q150 --> q151
    q151 -->|"tok(RegexLiteral)"| q152
    q152 --> q153
    q153 -->|"tok(&quot;;&quot;)"| q154
    q154 --> q19
```

## Alternatives

```mermaid
flowchart TD
    q20(["SN:20<br/>RuleStart"])
    q21(["SN:21<br/>RuleStop"])
    q155["SN:155<br/>Basic<br/>"]
    q156["SN:156<br/>Basic<br/>"]
    q157["SN:157<br/>Basic<br/>"]
    q158["SN:158<br/>Basic<br/>"]
    q159["SN:159<br/>Basic<br/>"]
    q160["SN:160<br/>Basic<br/>"]
    q161{"SN:161<br/>PlusLoopBack<br/><br/>dec=3"}
    q162["SN:162<br/>LoopEnd<br/>"]

    q20 --> q155
    q155 -.->|"[Group]"| q156
    q156 --> q157
    q157 -->|"tok(&quot;|&quot;)"| q158
    q157 --> q162
    q158 --> q159
    q159 -.->|"[Group]"| q160
    q160 --> q161
    q161 --> q157
    q161 --> q162
    q162 --> q21
```

## Group

```mermaid
flowchart TD
    q22(["SN:22<br/>RuleStart"])
    q23(["SN:23<br/>RuleStop"])
    q163["SN:163<br/>Basic<br/>"]
    q164["SN:164<br/>Basic<br/>"]
    q165["SN:165<br/>Basic<br/>"]
    q166["SN:166<br/>Basic<br/>"]
    q167{"SN:167<br/>PlusLoopBack<br/><br/>dec=4"}
    q168["SN:168<br/>LoopEnd<br/>"]

    q22 --> q163
    q163 -.->|"[Element]"| q164
    q164 --> q165
    q165 -.->|"[Element]"| q166
    q165 --> q168
    q166 --> q167
    q167 --> q165
    q167 --> q168
    q168 --> q23
```

## Element

```mermaid
flowchart TD
    q24(["SN:24<br/>RuleStart"])
    q25(["SN:25<br/>RuleStop"])
    q169["SN:169<br/>Basic<br/>"]
    q170["SN:170<br/>Basic<br/>"]
    q171["SN:171<br/>Basic<br/>"]
    q172["SN:172<br/>Basic<br/>"]
    q173["SN:173<br/>Basic<br/>"]
    q174["SN:174<br/>Basic<br/>"]
    q175["SN:175<br/>Basic<br/>"]
    q176["SN:176<br/>Basic<br/>"]
    q177["SN:177<br/>Basic<br/>"]
    q178["SN:178<br/>Basic<br/>"]
    q179["SN:179<br/>Basic<br/>"]
    q180["SN:180<br/>Basic<br/>"]
    q181["SN:181<br/>Basic<br/>"]
    q182["SN:182<br/>Basic<br/>"]
    q183["SN:183<br/>Basic<br/>"]
    q184["SN:184<br/>BlockEnd<br/>"]
    q185["SN:185<br/>Basic<br/>"]
    q186["SN:186<br/>Basic<br/>"]
    q187["SN:187<br/>Basic<br/>"]
    q188["SN:188<br/>Basic<br/>"]
    q189["SN:189<br/>Basic<br/>"]
    q190["SN:190<br/>Basic<br/>"]
    q191["SN:191<br/>Basic<br/>"]
    q192["SN:192<br/>BlockEnd<br/>"]

    q24 --> q183
    q169 -.->|"[Keyword]"| q170
    q170 --> q184
    q171 -.->|"[Assignment]"| q172
    q172 --> q184
    q173 -.->|"[RuleCall]"| q174
    q174 --> q184
    q175 -.->|"[Action]"| q176
    q176 --> q184
    q177 -->|"tok(&quot;(&quot;)"| q178
    q178 --> q179
    q179 -.->|"[Alternatives]"| q180
    q180 --> q181
    q181 -->|"tok(&quot;)&quot;)"| q182
    q182 --> q184
    q183 --> q169
    q183 --> q171
    q183 --> q173
    q183 --> q175
    q183 --> q177
    q184 --> q191
    q185 -->|"tok(&quot;*&quot;)"| q186
    q186 --> q192
    q187 -->|"tok(&quot;+&quot;)"| q188
    q188 --> q192
    q189 -->|"tok(&quot;?&quot;)"| q190
    q190 --> q192
    q191 --> q185
    q191 --> q187
    q191 --> q189
    q191 --> q192
    q192 --> q25
```

## Keyword

```mermaid
flowchart TD
    q26(["SN:26<br/>RuleStart"])
    q27(["SN:27<br/>RuleStop"])
    q193["SN:193<br/>Basic<br/>"]
    q194["SN:194<br/>Basic<br/>"]

    q26 --> q193
    q193 -->|"tok(StringLiteral)"| q194
    q194 --> q27
```

## Assignment

```mermaid
flowchart TD
    q28(["SN:28<br/>RuleStart"])
    q29(["SN:29<br/>RuleStop"])
    q195["SN:195<br/>Basic<br/>"]
    q196["SN:196<br/>Basic<br/>"]
    q197["SN:197<br/>Basic<br/>"]
    q198["SN:198<br/>Basic<br/>"]
    q199["SN:199<br/>Basic<br/>"]
    q200["SN:200<br/>Basic<br/>"]
    q201["SN:201<br/>Basic<br/>"]
    q202["SN:202<br/>Basic<br/>"]
    q203["SN:203<br/>Basic<br/>"]
    q204["SN:204<br/>BlockEnd<br/>"]
    q205["SN:205<br/>Basic<br/>"]
    q206["SN:206<br/>Basic<br/>"]

    q28 --> q195
    q195 -->|"tok(ID)"| q196
    q196 --> q203
    q197 -->|"tok(&quot;+=&quot;)"| q198
    q198 --> q204
    q199 -->|"tok(&quot;=&quot;)"| q200
    q200 --> q204
    q201 -->|"tok(&quot;?=&quot;)"| q202
    q202 --> q204
    q203 --> q197
    q203 --> q199
    q203 --> q201
    q204 --> q205
    q205 -.->|"[Assignable]"| q206
    q206 --> q29
```

## Assignable

```mermaid
flowchart TD
    q30(["SN:30<br/>RuleStart"])
    q31(["SN:31<br/>RuleStop"])
    q207["SN:207<br/>Basic<br/>"]
    q208["SN:208<br/>Basic<br/>"]
    q209["SN:209<br/>Basic<br/>"]
    q210["SN:210<br/>Basic<br/>"]
    q211["SN:211<br/>Basic<br/>"]
    q212["SN:212<br/>Basic<br/>"]
    q213["SN:213<br/>Basic<br/>"]
    q214["SN:214<br/>Basic<br/>"]
    q215["SN:215<br/>Basic<br/>"]
    q216["SN:216<br/>Basic<br/>"]
    q217["SN:217<br/>Basic<br/>"]
    q218["SN:218<br/>Basic<br/>"]
    q219["SN:219<br/>Basic<br/>"]
    q220["SN:220<br/>BlockEnd<br/>"]

    q30 --> q219
    q207 -.->|"[Keyword]"| q208
    q208 --> q220
    q209 -.->|"[RuleCall]"| q210
    q210 --> q220
    q211 -.->|"[CrossRef]"| q212
    q212 --> q220
    q213 -->|"tok(&quot;(&quot;)"| q214
    q214 --> q215
    q215 -.->|"[AssignableAlternatives]"| q216
    q216 --> q217
    q217 -->|"tok(&quot;)&quot;)"| q218
    q218 --> q220
    q219 --> q207
    q219 --> q209
    q219 --> q211
    q219 --> q213
    q220 --> q31
```

## AssignableWithoutAlts

```mermaid
flowchart TD
    q32(["SN:32<br/>RuleStart"])
    q33(["SN:33<br/>RuleStop"])
    q221["SN:221<br/>Basic<br/>"]
    q222["SN:222<br/>Basic<br/>"]
    q223["SN:223<br/>Basic<br/>"]
    q224["SN:224<br/>Basic<br/>"]
    q225["SN:225<br/>Basic<br/>"]
    q226["SN:226<br/>Basic<br/>"]
    q227["SN:227<br/>Basic<br/>"]
    q228["SN:228<br/>BlockEnd<br/>"]

    q32 --> q227
    q221 -.->|"[Keyword]"| q222
    q222 --> q228
    q223 -.->|"[RuleCall]"| q224
    q224 --> q228
    q225 -.->|"[CrossRef]"| q226
    q226 --> q228
    q227 --> q221
    q227 --> q223
    q227 --> q225
    q228 --> q33
```

## AssignableAlternatives

```mermaid
flowchart TD
    q34(["SN:34<br/>RuleStart"])
    q35(["SN:35<br/>RuleStop"])
    q229["SN:229<br/>Basic<br/>"]
    q230["SN:230<br/>Basic<br/>"]
    q231["SN:231<br/>Basic<br/>"]
    q232["SN:232<br/>Basic<br/>"]
    q233["SN:233<br/>Basic<br/>"]
    q234["SN:234<br/>Basic<br/>"]
    q235{"SN:235<br/>PlusLoopBack<br/><br/>dec=5"}
    q236["SN:236<br/>LoopEnd<br/>"]

    q34 --> q229
    q229 -.->|"[AssignableWithoutAlts]"| q230
    q230 --> q231
    q231 -->|"tok(&quot;|&quot;)"| q232
    q231 --> q236
    q232 --> q233
    q233 -.->|"[AssignableWithoutAlts]"| q234
    q234 --> q235
    q235 --> q231
    q235 --> q236
    q236 --> q35
```

## CrossRef

```mermaid
flowchart TD
    q36(["SN:36<br/>RuleStart"])
    q37(["SN:37<br/>RuleStop"])
    q237["SN:237<br/>Basic<br/>"]
    q238["SN:238<br/>Basic<br/>"]
    q239["SN:239<br/>Basic<br/>"]
    q240["SN:240<br/>Basic<br/>"]
    q241["SN:241<br/>Basic<br/>"]
    q242["SN:242<br/>Basic<br/>"]
    q243["SN:243<br/>Basic<br/>"]
    q244["SN:244<br/>Basic<br/>"]
    q245["SN:245<br/>Basic<br/>"]
    q246["SN:246<br/>Basic<br/>"]

    q36 --> q237
    q237 -->|"tok(&quot;[&quot;)"| q238
    q238 --> q239
    q239 -->|"tok(ID)"| q240
    q240 --> q241
    q241 -->|"tok(&quot;:&quot;)"| q242
    q241 --> q244
    q242 --> q243
    q243 -.->|"[RuleCall]"| q244
    q244 --> q245
    q245 -->|"tok(&quot;]&quot;)"| q246
    q246 --> q37
```

## RuleCall

```mermaid
flowchart TD
    q38(["SN:38<br/>RuleStart"])
    q39(["SN:39<br/>RuleStop"])
    q247["SN:247<br/>Basic<br/>"]
    q248["SN:248<br/>Basic<br/>"]

    q38 --> q247
    q247 -->|"tok(ID)"| q248
    q248 --> q39
```

## Action

```mermaid
flowchart TD
    q40(["SN:40<br/>RuleStart"])
    q41(["SN:41<br/>RuleStop"])
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
    q262["SN:262<br/>BlockEnd<br/>"]
    q263["SN:263<br/>Basic<br/>"]
    q264["SN:264<br/>Basic<br/>"]
    q265["SN:265<br/>Basic<br/>"]
    q266["SN:266<br/>Basic<br/>"]

    q40 --> q249
    q249 -->|"tok(&quot;{&quot;)"| q250
    q250 --> q251
    q251 -->|"tok(ID)"| q252
    q252 --> q253
    q253 -->|"tok(&quot;.&quot;)"| q254
    q253 --> q264
    q254 --> q255
    q255 -->|"tok(ID)"| q256
    q256 --> q261
    q257 -->|"tok(&quot;+=&quot;)"| q258
    q258 --> q262
    q259 -->|"tok(&quot;=&quot;)"| q260
    q260 --> q262
    q261 --> q257
    q261 --> q259
    q262 --> q263
    q263 -->|"tok(&quot;current&quot;)"| q264
    q264 --> q265
    q265 -->|"tok(&quot;}&quot;)"| q266
    q266 --> q41
```

## CompositeRule

```mermaid
flowchart TD
    q42(["SN:42<br/>RuleStart"])
    q43(["SN:43<br/>RuleStop"])
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

    q42 --> q267
    q267 -->|"tok(&quot;composite&quot;)"| q268
    q268 --> q269
    q269 -->|"tok(ID)"| q270
    q270 --> q271
    q271 -->|"tok(&quot;:&quot;)"| q272
    q272 --> q273
    q273 -.->|"[CompositeAlternatives]"| q274
    q274 --> q275
    q275 -->|"tok(&quot;;&quot;)"| q276
    q276 --> q43
```

## CompositeAlternatives

```mermaid
flowchart TD
    q44(["SN:44<br/>RuleStart"])
    q45(["SN:45<br/>RuleStop"])
    q277["SN:277<br/>Basic<br/>"]
    q278["SN:278<br/>Basic<br/>"]
    q279["SN:279<br/>Basic<br/>"]
    q280["SN:280<br/>Basic<br/>"]
    q281["SN:281<br/>Basic<br/>"]
    q282["SN:282<br/>Basic<br/>"]
    q283{"SN:283<br/>PlusLoopBack<br/><br/>dec=6"}
    q284["SN:284<br/>LoopEnd<br/>"]

    q44 --> q277
    q277 -.->|"[CompositeGroup]"| q278
    q278 --> q279
    q279 -->|"tok(&quot;|&quot;)"| q280
    q279 --> q284
    q280 --> q281
    q281 -.->|"[CompositeGroup]"| q282
    q282 --> q283
    q283 --> q279
    q283 --> q284
    q284 --> q45
```

## CompositeGroup

```mermaid
flowchart TD
    q46(["SN:46<br/>RuleStart"])
    q47(["SN:47<br/>RuleStop"])
    q285["SN:285<br/>Basic<br/>"]
    q286["SN:286<br/>Basic<br/>"]
    q287["SN:287<br/>Basic<br/>"]
    q288["SN:288<br/>Basic<br/>"]
    q289{"SN:289<br/>PlusLoopBack<br/><br/>dec=7"}
    q290["SN:290<br/>LoopEnd<br/>"]

    q46 --> q285
    q285 -.->|"[CompositeElement]"| q286
    q286 --> q287
    q287 -.->|"[CompositeElement]"| q288
    q287 --> q290
    q288 --> q289
    q289 --> q287
    q289 --> q290
    q290 --> q47
```

## CompositeElement

```mermaid
flowchart TD
    q48(["SN:48<br/>RuleStart"])
    q49(["SN:49<br/>RuleStop"])
    q291["SN:291<br/>Basic<br/>"]
    q292["SN:292<br/>Basic<br/>"]
    q293["SN:293<br/>Basic<br/>"]
    q294["SN:294<br/>Basic<br/>"]
    q295["SN:295<br/>Basic<br/>"]
    q296["SN:296<br/>Basic<br/>"]
    q297["SN:297<br/>Basic<br/>"]
    q298["SN:298<br/>Basic<br/>"]
    q299["SN:299<br/>Basic<br/>"]
    q300["SN:300<br/>Basic<br/>"]
    q301["SN:301<br/>Basic<br/>"]
    q302["SN:302<br/>BlockEnd<br/>"]
    q303["SN:303<br/>Basic<br/>"]
    q304["SN:304<br/>Basic<br/>"]
    q305["SN:305<br/>Basic<br/>"]
    q306["SN:306<br/>Basic<br/>"]
    q307["SN:307<br/>Basic<br/>"]
    q308["SN:308<br/>Basic<br/>"]
    q309["SN:309<br/>Basic<br/>"]
    q310["SN:310<br/>BlockEnd<br/>"]

    q48 --> q301
    q291 -.->|"[Keyword]"| q292
    q292 --> q302
    q293 -.->|"[RuleCall]"| q294
    q294 --> q302
    q295 -->|"tok(&quot;(&quot;)"| q296
    q296 --> q297
    q297 -.->|"[CompositeAlternatives]"| q298
    q298 --> q299
    q299 -->|"tok(&quot;)&quot;)"| q300
    q300 --> q302
    q301 --> q291
    q301 --> q293
    q301 --> q295
    q302 --> q309
    q303 -->|"tok(&quot;*&quot;)"| q304
    q304 --> q310
    q305 -->|"tok(&quot;+&quot;)"| q306
    q306 --> q310
    q307 -->|"tok(&quot;?&quot;)"| q308
    q308 --> q310
    q309 --> q303
    q309 --> q305
    q309 --> q307
    q309 --> q310
    q310 --> q49
```

