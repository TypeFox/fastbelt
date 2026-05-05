# Runtime ATN for grammar

## Grammar

```mermaid
flowchart TD
    q0(["SN:0<br/>RuleStart"])
    q1(["SN:1<br/>RuleStop"])
    q282["SN:282<br/>Basic<br/> #0"]
    q283["SN:283<br/>Basic<br/> #0"]
    q284["SN:284<br/>Basic<br/> #0"]
    q285["SN:285<br/>Basic<br/> #0"]
    q286["SN:286<br/>Basic<br/> #0"]
    q287["SN:287<br/>Basic<br/> #0"]
    q288["SN:288<br/>Basic<br/> #0"]
    q289["SN:289<br/>Basic<br/> #0"]
    q290["SN:290<br/>Basic<br/> #0"]
    q291["SN:291<br/>Basic<br/> #0"]
    q292["SN:292<br/>Basic<br/> #0"]
    q293["SN:293<br/>Basic<br/> #0"]
    q294["SN:294<br/>Basic<br/> #0"]
    q295["SN:295<br/>Basic<br/> #0"]
    q296["SN:296<br/>Basic<br/> #0"]
    q297["SN:297<br/>BlockEnd<br/> #0"]

    q0 --> q282
    q282 -->|"tok("grammar")"| q283
    q283 --> q284
    q284 -->|"tok(ID)"| q285
    q285 --> q286
    q286 -->|"tok(";")"| q287
    q287 --> q296
    q288 -.->|"[ParserRule]"| q289
    q289 --> q297
    q290 -.->|"[Token]"| q291
    q291 --> q297
    q292 -.->|"[Interface]"| q293
    q293 --> q297
    q294 -.->|"[CompositeRule]"| q295
    q295 --> q297
    q296 --> q288
    q296 --> q290
    q296 --> q292
    q296 --> q294
    q297 --> q1
```

## Interface

```mermaid
flowchart TD
    q2(["SN:2<br/>RuleStart"])
    q3(["SN:3<br/>RuleStop"])
    q118["SN:118<br/>Basic<br/> #0"]
    q119["SN:119<br/>Basic<br/> #0"]
    q120["SN:120<br/>Basic<br/> #0"]
    q121["SN:121<br/>Basic<br/> #0"]
    q122["SN:122<br/>Basic<br/> #0"]
    q123["SN:123<br/>Basic<br/> #0"]
    q124["SN:124<br/>Basic<br/> #0"]
    q125["SN:125<br/>Basic<br/> #0"]
    q126["SN:126<br/>Basic<br/> #0"]
    q127["SN:127<br/>Basic<br/> #0"]
    q128["SN:128<br/>Basic<br/> #0"]
    q129["SN:129<br/>Basic<br/> #0"]
    q130{"SN:130<br/>StarLoopEntry<br/> #0<br/>dec=1"}
    q131["SN:131<br/>LoopEnd<br/> #0"]
    q132["SN:132<br/>StarLoopBack<br/> #0"]
    q133["SN:133<br/>Basic<br/> #0"]
    q134["SN:134<br/>Basic<br/> #0"]
    q135["SN:135<br/>Basic<br/> #0"]
    q136["SN:136<br/>Basic<br/> #0"]
    q137{"SN:137<br/>StarLoopEntry<br/> #0<br/>dec=2"}
    q138["SN:138<br/>LoopEnd<br/> #0"]
    q139["SN:139<br/>StarLoopBack<br/> #0"]
    q140["SN:140<br/>Basic<br/> #0"]
    q141["SN:141<br/>Basic<br/> #0"]

    q2 --> q118
    q118 -->|"tok("interface")"| q119
    q119 --> q120
    q120 -->|"tok(ID)"| q121
    q121 --> q122
    q122 -->|"tok("extends")"| q123
    q122 --> q131
    q123 --> q124
    q124 -->|"tok(ID)"| q125
    q125 --> q130
    q126 -->|"tok(",")"| q127
    q127 --> q128
    q128 -->|"tok(ID)"| q129
    q129 --> q132
    q130 --> q126
    q130 --> q131
    q131 --> q133
    q132 --> q130
    q133 -->|"tok("{")"| q134
    q134 --> q137
    q135 -.->|"[Field]"| q136
    q136 --> q139
    q137 --> q135
    q137 --> q138
    q138 --> q140
    q139 --> q137
    q140 -->|"tok("}")"| q141
    q141 --> q3
```

## Field

```mermaid
flowchart TD
    q4(["SN:4<br/>RuleStart"])
    q5(["SN:5<br/>RuleStop"])
    q142["SN:142<br/>Basic<br/> #0"]
    q143["SN:143<br/>Basic<br/> #0"]
    q144["SN:144<br/>Basic<br/> #0"]
    q145["SN:145<br/>Basic<br/> #0"]

    q4 --> q142
    q142 -->|"tok(ID)"| q143
    q143 --> q144
    q144 -.->|"[FieldType]"| q145
    q145 --> q5
```

## FieldType

```mermaid
flowchart TD
    q6(["SN:6<br/>RuleStart"])
    q7(["SN:7<br/>RuleStop"])
    q298["SN:298<br/>Basic<br/> #0"]
    q299["SN:299<br/>Basic<br/> #0"]
    q300["SN:300<br/>Basic<br/> #0"]
    q301["SN:301<br/>Basic<br/> #0"]
    q302["SN:302<br/>Basic<br/> #0"]
    q303["SN:303<br/>Basic<br/> #0"]
    q304["SN:304<br/>Basic<br/> #0"]
    q305["SN:305<br/>Basic<br/> #0"]
    q306["SN:306<br/>Basic<br/> #0"]
    q307["SN:307<br/>BlockEnd<br/> #0"]

    q6 --> q306
    q298 -.->|"[SimpleType]"| q299
    q299 --> q307
    q300 -.->|"[ReferenceType]"| q301
    q301 --> q307
    q302 -.->|"[ArrayType]"| q303
    q303 --> q307
    q304 -.->|"[PrimitiveType]"| q305
    q305 --> q307
    q306 --> q298
    q306 --> q300
    q306 --> q302
    q306 --> q304
    q307 --> q7
```

## ArrayType

```mermaid
flowchart TD
    q8(["SN:8<br/>RuleStart"])
    q9(["SN:9<br/>RuleStop"])
    q182["SN:182<br/>Basic<br/> #0"]
    q183["SN:183<br/>Basic<br/> #0"]
    q184["SN:184<br/>Basic<br/> #0"]
    q185["SN:185<br/>Basic<br/> #0"]
    q186["SN:186<br/>Basic<br/> #0"]
    q187["SN:187<br/>Basic<br/> #0"]

    q8 --> q182
    q182 -->|"tok("[")"| q183
    q183 --> q184
    q184 -->|"tok("]")"| q185
    q185 --> q186
    q186 -.->|"[FieldType]"| q187
    q187 --> q9
```

## ReferenceType

```mermaid
flowchart TD
    q10(["SN:10<br/>RuleStart"])
    q11(["SN:11<br/>RuleStop"])
    q188["SN:188<br/>Basic<br/> #0"]
    q189["SN:189<br/>Basic<br/> #0"]
    q190["SN:190<br/>Basic<br/> #0"]
    q191["SN:191<br/>Basic<br/> #0"]

    q10 --> q188
    q188 -->|"tok("*")"| q189
    q189 --> q190
    q190 -->|"tok(ID)"| q191
    q191 --> q11
```

## SimpleType

```mermaid
flowchart TD
    q12(["SN:12<br/>RuleStart"])
    q13(["SN:13<br/>RuleStop"])
    q146["SN:146<br/>Basic<br/> #0"]
    q147["SN:147<br/>Basic<br/> #0"]

    q12 --> q146
    q146 -->|"tok(ID)"| q147
    q147 --> q13
```

## PrimitiveType

```mermaid
flowchart TD
    q14(["SN:14<br/>RuleStart"])
    q15(["SN:15<br/>RuleStop"])
    q148["SN:148<br/>Basic<br/> #0"]
    q149["SN:149<br/>Basic<br/> #0"]
    q150["SN:150<br/>Basic<br/> #0"]
    q151["SN:151<br/>Basic<br/> #0"]
    q152["SN:152<br/>Basic<br/> #0"]
    q153["SN:153<br/>Basic<br/> #0"]
    q154["SN:154<br/>Basic<br/> #0"]
    q155["SN:155<br/>BlockEnd<br/> #0"]

    q14 --> q154
    q148 -->|"tok("string")"| q149
    q149 --> q155
    q150 -->|"tok("bool")"| q151
    q151 --> q155
    q152 -->|"tok("composite")"| q153
    q153 --> q155
    q154 --> q148
    q154 --> q150
    q154 --> q152
    q155 --> q15
```

## ParserRule

```mermaid
flowchart TD
    q16(["SN:16<br/>RuleStart"])
    q17(["SN:17<br/>RuleStop"])
    q192["SN:192<br/>Basic<br/> #0"]
    q193["SN:193<br/>Basic<br/> #0"]
    q194["SN:194<br/>Basic<br/> #0"]
    q195["SN:195<br/>Basic<br/> #0"]
    q196["SN:196<br/>Basic<br/> #0"]
    q197["SN:197<br/>Basic<br/> #0"]
    q198["SN:198<br/>Basic<br/> #0"]
    q199["SN:199<br/>Basic<br/> #0"]
    q200["SN:200<br/>Basic<br/> #0"]
    q201["SN:201<br/>Basic<br/> #0"]
    q202["SN:202<br/>Basic<br/> #0"]
    q203["SN:203<br/>Basic<br/> #0"]

    q16 --> q192
    q192 -->|"tok(ID)"| q193
    q193 --> q194
    q194 -->|"tok("returns")"| q195
    q194 --> q197
    q195 --> q196
    q196 -->|"tok(ID)"| q197
    q197 --> q198
    q198 -->|"tok(":")"| q199
    q199 --> q200
    q200 -.->|"[Alternatives]"| q201
    q201 --> q202
    q202 -->|"tok(";")"| q203
    q203 --> q17
```

## Token

```mermaid
flowchart TD
    q18(["SN:18<br/>RuleStart"])
    q19(["SN:19<br/>RuleStop"])
    q50["SN:50<br/>Basic<br/> #0"]
    q51["SN:51<br/>Basic<br/> #0"]
    q52["SN:52<br/>Basic<br/> #0"]
    q53["SN:53<br/>Basic<br/> #0"]
    q54["SN:54<br/>Basic<br/> #0"]
    q55["SN:55<br/>BlockEnd<br/> #0"]
    q56["SN:56<br/>Basic<br/> #0"]
    q57["SN:57<br/>Basic<br/> #0"]
    q58["SN:58<br/>Basic<br/> #0"]
    q59["SN:59<br/>Basic<br/> #0"]
    q60["SN:60<br/>Basic<br/> #0"]
    q61["SN:61<br/>Basic<br/> #0"]
    q62["SN:62<br/>Basic<br/> #0"]
    q63["SN:63<br/>Basic<br/> #0"]
    q64["SN:64<br/>Basic<br/> #0"]
    q65["SN:65<br/>Basic<br/> #0"]

    q18 --> q54
    q50 -->|"tok("hidden")"| q51
    q51 --> q55
    q52 -->|"tok("comment")"| q53
    q53 --> q55
    q54 --> q50
    q54 --> q52
    q55 --> q56
    q56 -->|"tok("token")"| q57
    q57 --> q58
    q58 -->|"tok(ID)"| q59
    q59 --> q60
    q60 -->|"tok(":")"| q61
    q61 --> q62
    q62 -->|"tok(RegexLiteral)"| q63
    q63 --> q64
    q64 -->|"tok(";")"| q65
    q65 --> q19
```

## Alternatives

```mermaid
flowchart TD
    q20(["SN:20<br/>RuleStart"])
    q21(["SN:21<br/>RuleStop"])
    q66["SN:66<br/>Basic<br/> #0"]
    q67["SN:67<br/>Basic<br/> #0"]
    q68["SN:68<br/>Basic<br/> #0"]
    q69["SN:69<br/>Basic<br/> #0"]
    q70["SN:70<br/>Basic<br/> #0"]
    q71["SN:71<br/>Basic<br/> #0"]
    q72{"SN:72<br/>PlusLoopBack<br/> #0<br/>dec=0"}
    q73["SN:73<br/>LoopEnd<br/> #0"]

    q20 --> q66
    q66 -.->|"[Group]"| q67
    q67 --> q68
    q68 -->|"tok("|")"| q69
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
    q156["SN:156<br/>Basic<br/> #0"]
    q157["SN:157<br/>Basic<br/> #0"]
    q158["SN:158<br/>Basic<br/> #0"]
    q159["SN:159<br/>Basic<br/> #0"]
    q160{"SN:160<br/>PlusLoopBack<br/> #0<br/>dec=3"}
    q161["SN:161<br/>LoopEnd<br/> #0"]

    q22 --> q156
    q156 -.->|"[Element]"| q157
    q157 --> q158
    q158 -.->|"[Element]"| q159
    q158 --> q161
    q159 --> q160
    q160 --> q158
    q160 --> q161
    q161 --> q23
```

## Element

```mermaid
flowchart TD
    q24(["SN:24<br/>RuleStart"])
    q25(["SN:25<br/>RuleStop"])
    q74["SN:74<br/>Basic<br/> #0"]
    q75["SN:75<br/>Basic<br/> #0"]
    q76["SN:76<br/>Basic<br/> #0"]
    q77["SN:77<br/>Basic<br/> #0"]
    q78["SN:78<br/>Basic<br/> #0"]
    q79["SN:79<br/>Basic<br/> #0"]
    q80["SN:80<br/>Basic<br/> #0"]
    q81["SN:81<br/>Basic<br/> #0"]
    q82["SN:82<br/>Basic<br/> #0"]
    q83["SN:83<br/>Basic<br/> #0"]
    q84["SN:84<br/>Basic<br/> #0"]
    q85["SN:85<br/>Basic<br/> #0"]
    q86["SN:86<br/>Basic<br/> #0"]
    q87["SN:87<br/>Basic<br/> #0"]
    q88["SN:88<br/>Basic<br/> #0"]
    q89["SN:89<br/>BlockEnd<br/> #0"]
    q90["SN:90<br/>Basic<br/> #0"]
    q91["SN:91<br/>Basic<br/> #0"]
    q92["SN:92<br/>Basic<br/> #0"]
    q93["SN:93<br/>Basic<br/> #0"]
    q94["SN:94<br/>Basic<br/> #0"]
    q95["SN:95<br/>Basic<br/> #0"]
    q96["SN:96<br/>Basic<br/> #0"]
    q97["SN:97<br/>BlockEnd<br/> #0"]

    q24 --> q88
    q74 -.->|"[Keyword]"| q75
    q75 --> q89
    q76 -.->|"[Assignment]"| q77
    q77 --> q89
    q78 -.->|"[RuleCall]"| q79
    q79 --> q89
    q80 -.->|"[Action]"| q81
    q81 --> q89
    q82 -->|"tok("(")"| q83
    q83 --> q84
    q84 -.->|"[Alternatives]"| q85
    q85 --> q86
    q86 -->|"tok(")")"| q87
    q87 --> q89
    q88 --> q74
    q88 --> q76
    q88 --> q78
    q88 --> q80
    q88 --> q82
    q89 --> q96
    q90 -->|"tok("*")"| q91
    q91 --> q97
    q92 -->|"tok("+")"| q93
    q93 --> q97
    q94 -->|"tok("?")"| q95
    q95 --> q97
    q96 --> q90
    q96 --> q92
    q96 --> q94
    q97 --> q25
```

## Keyword

```mermaid
flowchart TD
    q26(["SN:26<br/>RuleStart"])
    q27(["SN:27<br/>RuleStop"])
    q238["SN:238<br/>Basic<br/> #0"]
    q239["SN:239<br/>Basic<br/> #0"]

    q26 --> q238
    q238 -->|"tok(StringLiteral)"| q239
    q239 --> q27
```

## Assignment

```mermaid
flowchart TD
    q28(["SN:28<br/>RuleStart"])
    q29(["SN:29<br/>RuleStop"])
    q98["SN:98<br/>Basic<br/> #0"]
    q99["SN:99<br/>Basic<br/> #0"]
    q100["SN:100<br/>Basic<br/> #0"]
    q101["SN:101<br/>Basic<br/> #0"]
    q102["SN:102<br/>Basic<br/> #0"]
    q103["SN:103<br/>Basic<br/> #0"]
    q104["SN:104<br/>Basic<br/> #0"]
    q105["SN:105<br/>Basic<br/> #0"]
    q106["SN:106<br/>Basic<br/> #0"]
    q107["SN:107<br/>BlockEnd<br/> #0"]
    q108["SN:108<br/>Basic<br/> #0"]
    q109["SN:109<br/>Basic<br/> #0"]

    q28 --> q98
    q98 -->|"tok(ID)"| q99
    q99 --> q106
    q100 -->|"tok("+=")"| q101
    q101 --> q107
    q102 -->|"tok("=")"| q103
    q103 --> q107
    q104 -->|"tok("?=")"| q105
    q105 --> q107
    q106 --> q100
    q106 --> q102
    q106 --> q104
    q107 --> q108
    q108 -.->|"[Assignable]"| q109
    q109 --> q29
```

## Assignable

```mermaid
flowchart TD
    q30(["SN:30<br/>RuleStart"])
    q31(["SN:31<br/>RuleStop"])
    q204["SN:204<br/>Basic<br/> #0"]
    q205["SN:205<br/>Basic<br/> #0"]
    q206["SN:206<br/>Basic<br/> #0"]
    q207["SN:207<br/>Basic<br/> #0"]
    q208["SN:208<br/>Basic<br/> #0"]
    q209["SN:209<br/>Basic<br/> #0"]
    q210["SN:210<br/>Basic<br/> #0"]
    q211["SN:211<br/>Basic<br/> #0"]
    q212["SN:212<br/>Basic<br/> #0"]
    q213["SN:213<br/>Basic<br/> #0"]
    q214["SN:214<br/>Basic<br/> #0"]
    q215["SN:215<br/>Basic<br/> #0"]
    q216["SN:216<br/>Basic<br/> #0"]
    q217["SN:217<br/>BlockEnd<br/> #0"]

    q30 --> q216
    q204 -.->|"[Keyword]"| q205
    q205 --> q217
    q206 -.->|"[RuleCall]"| q207
    q207 --> q217
    q208 -.->|"[CrossRef]"| q209
    q209 --> q217
    q210 -->|"tok("(")"| q211
    q211 --> q212
    q212 -.->|"[AssignableAlternatives]"| q213
    q213 --> q214
    q214 -->|"tok(")")"| q215
    q215 --> q217
    q216 --> q204
    q216 --> q206
    q216 --> q208
    q216 --> q210
    q217 --> q31
```

## AssignableWithoutAlts

```mermaid
flowchart TD
    q32(["SN:32<br/>RuleStart"])
    q33(["SN:33<br/>RuleStop"])
    q110["SN:110<br/>Basic<br/> #0"]
    q111["SN:111<br/>Basic<br/> #0"]
    q112["SN:112<br/>Basic<br/> #0"]
    q113["SN:113<br/>Basic<br/> #0"]
    q114["SN:114<br/>Basic<br/> #0"]
    q115["SN:115<br/>Basic<br/> #0"]
    q116["SN:116<br/>Basic<br/> #0"]
    q117["SN:117<br/>BlockEnd<br/> #0"]

    q32 --> q116
    q110 -.->|"[Keyword]"| q111
    q111 --> q117
    q112 -.->|"[RuleCall]"| q113
    q113 --> q117
    q114 -.->|"[CrossRef]"| q115
    q115 --> q117
    q116 --> q110
    q116 --> q112
    q116 --> q114
    q117 --> q33
```

## AssignableAlternatives

```mermaid
flowchart TD
    q34(["SN:34<br/>RuleStart"])
    q35(["SN:35<br/>RuleStop"])
    q162["SN:162<br/>Basic<br/> #0"]
    q163["SN:163<br/>Basic<br/> #0"]
    q164["SN:164<br/>Basic<br/> #0"]
    q165["SN:165<br/>Basic<br/> #0"]
    q166["SN:166<br/>Basic<br/> #0"]
    q167["SN:167<br/>Basic<br/> #0"]
    q168{"SN:168<br/>PlusLoopBack<br/> #0<br/>dec=4"}
    q169["SN:169<br/>LoopEnd<br/> #0"]

    q34 --> q162
    q162 -.->|"[AssignableWithoutAlts]"| q163
    q163 --> q164
    q164 -->|"tok("|")"| q165
    q164 --> q169
    q165 --> q166
    q166 -.->|"[AssignableWithoutAlts]"| q167
    q167 --> q168
    q168 --> q164
    q168 --> q169
    q169 --> q35
```

## CrossRef

```mermaid
flowchart TD
    q36(["SN:36<br/>RuleStart"])
    q37(["SN:37<br/>RuleStop"])
    q240["SN:240<br/>Basic<br/> #0"]
    q241["SN:241<br/>Basic<br/> #0"]
    q242["SN:242<br/>Basic<br/> #0"]
    q243["SN:243<br/>Basic<br/> #0"]
    q244["SN:244<br/>Basic<br/> #0"]
    q245["SN:245<br/>Basic<br/> #0"]
    q246["SN:246<br/>Basic<br/> #0"]
    q247["SN:247<br/>Basic<br/> #0"]
    q248["SN:248<br/>Basic<br/> #0"]
    q249["SN:249<br/>Basic<br/> #0"]

    q36 --> q240
    q240 -->|"tok("[")"| q241
    q241 --> q242
    q242 -->|"tok(ID)"| q243
    q243 --> q244
    q244 -->|"tok(":")"| q245
    q244 --> q247
    q245 --> q246
    q246 -.->|"[RuleCall]"| q247
    q247 --> q248
    q248 -->|"tok("]")"| q249
    q249 --> q37
```

## RuleCall

```mermaid
flowchart TD
    q38(["SN:38<br/>RuleStart"])
    q39(["SN:39<br/>RuleStop"])
    q170["SN:170<br/>Basic<br/> #0"]
    q171["SN:171<br/>Basic<br/> #0"]

    q38 --> q170
    q170 -->|"tok(ID)"| q171
    q171 --> q39
```

## Action

```mermaid
flowchart TD
    q40(["SN:40<br/>RuleStart"])
    q41(["SN:41<br/>RuleStop"])
    q250["SN:250<br/>Basic<br/> #0"]
    q251["SN:251<br/>Basic<br/> #0"]
    q252["SN:252<br/>Basic<br/> #0"]
    q253["SN:253<br/>Basic<br/> #0"]
    q254["SN:254<br/>Basic<br/> #0"]
    q255["SN:255<br/>Basic<br/> #0"]
    q256["SN:256<br/>Basic<br/> #0"]
    q257["SN:257<br/>Basic<br/> #0"]
    q258["SN:258<br/>Basic<br/> #0"]
    q259["SN:259<br/>Basic<br/> #0"]
    q260["SN:260<br/>Basic<br/> #0"]
    q261["SN:261<br/>Basic<br/> #0"]
    q262["SN:262<br/>Basic<br/> #0"]
    q263["SN:263<br/>BlockEnd<br/> #0"]
    q264["SN:264<br/>Basic<br/> #0"]
    q265["SN:265<br/>Basic<br/> #0"]
    q266["SN:266<br/>Basic<br/> #0"]
    q267["SN:267<br/>Basic<br/> #0"]

    q40 --> q250
    q250 -->|"tok("{")"| q251
    q251 --> q252
    q252 -->|"tok(ID)"| q253
    q253 --> q254
    q254 -->|"tok(".")"| q255
    q254 --> q265
    q255 --> q256
    q256 -->|"tok(ID)"| q257
    q257 --> q262
    q258 -->|"tok("+=")"| q259
    q259 --> q263
    q260 -->|"tok("=")"| q261
    q261 --> q263
    q262 --> q258
    q262 --> q260
    q263 --> q264
    q264 -->|"tok("current")"| q265
    q265 --> q266
    q266 -->|"tok("}")"| q267
    q267 --> q41
```

## CompositeRule

```mermaid
flowchart TD
    q42(["SN:42<br/>RuleStart"])
    q43(["SN:43<br/>RuleStop"])
    q172["SN:172<br/>Basic<br/> #0"]
    q173["SN:173<br/>Basic<br/> #0"]
    q174["SN:174<br/>Basic<br/> #0"]
    q175["SN:175<br/>Basic<br/> #0"]
    q176["SN:176<br/>Basic<br/> #0"]
    q177["SN:177<br/>Basic<br/> #0"]
    q178["SN:178<br/>Basic<br/> #0"]
    q179["SN:179<br/>Basic<br/> #0"]
    q180["SN:180<br/>Basic<br/> #0"]
    q181["SN:181<br/>Basic<br/> #0"]

    q42 --> q172
    q172 -->|"tok("composite")"| q173
    q173 --> q174
    q174 -->|"tok(ID)"| q175
    q175 --> q176
    q176 -->|"tok(":")"| q177
    q177 --> q178
    q178 -.->|"[CompositeAlternatives]"| q179
    q179 --> q180
    q180 -->|"tok(";")"| q181
    q181 --> q43
```

## CompositeAlternatives

```mermaid
flowchart TD
    q44(["SN:44<br/>RuleStart"])
    q45(["SN:45<br/>RuleStop"])
    q268["SN:268<br/>Basic<br/> #0"]
    q269["SN:269<br/>Basic<br/> #0"]
    q270["SN:270<br/>Basic<br/> #0"]
    q271["SN:271<br/>Basic<br/> #0"]
    q272["SN:272<br/>Basic<br/> #0"]
    q273["SN:273<br/>Basic<br/> #0"]
    q274{"SN:274<br/>PlusLoopBack<br/> #0<br/>dec=5"}
    q275["SN:275<br/>LoopEnd<br/> #0"]

    q44 --> q268
    q268 -.->|"[CompositeGroup]"| q269
    q269 --> q270
    q270 -->|"tok("|")"| q271
    q270 --> q275
    q271 --> q272
    q272 -.->|"[CompositeGroup]"| q273
    q273 --> q274
    q274 --> q270
    q274 --> q275
    q275 --> q45
```

## CompositeGroup

```mermaid
flowchart TD
    q46(["SN:46<br/>RuleStart"])
    q47(["SN:47<br/>RuleStop"])
    q276["SN:276<br/>Basic<br/> #0"]
    q277["SN:277<br/>Basic<br/> #0"]
    q278["SN:278<br/>Basic<br/> #0"]
    q279["SN:279<br/>Basic<br/> #0"]
    q280{"SN:280<br/>PlusLoopBack<br/> #0<br/>dec=6"}
    q281["SN:281<br/>LoopEnd<br/> #0"]

    q46 --> q276
    q276 -.->|"[CompositeElement]"| q277
    q277 --> q278
    q278 -.->|"[CompositeElement]"| q279
    q278 --> q281
    q279 --> q280
    q280 --> q278
    q280 --> q281
    q281 --> q47
```

## CompositeElement

```mermaid
flowchart TD
    q48(["SN:48<br/>RuleStart"])
    q49(["SN:49<br/>RuleStop"])
    q218["SN:218<br/>Basic<br/> #0"]
    q219["SN:219<br/>Basic<br/> #0"]
    q220["SN:220<br/>Basic<br/> #0"]
    q221["SN:221<br/>Basic<br/> #0"]
    q222["SN:222<br/>Basic<br/> #0"]
    q223["SN:223<br/>Basic<br/> #0"]
    q224["SN:224<br/>Basic<br/> #0"]
    q225["SN:225<br/>Basic<br/> #0"]
    q226["SN:226<br/>Basic<br/> #0"]
    q227["SN:227<br/>Basic<br/> #0"]
    q228["SN:228<br/>Basic<br/> #0"]
    q229["SN:229<br/>BlockEnd<br/> #0"]
    q230["SN:230<br/>Basic<br/> #0"]
    q231["SN:231<br/>Basic<br/> #0"]
    q232["SN:232<br/>Basic<br/> #0"]
    q233["SN:233<br/>Basic<br/> #0"]
    q234["SN:234<br/>Basic<br/> #0"]
    q235["SN:235<br/>Basic<br/> #0"]
    q236["SN:236<br/>Basic<br/> #0"]
    q237["SN:237<br/>BlockEnd<br/> #0"]

    q48 --> q228
    q218 -.->|"[Keyword]"| q219
    q219 --> q229
    q220 -.->|"[RuleCall]"| q221
    q221 --> q229
    q222 -->|"tok("(")"| q223
    q223 --> q224
    q224 -.->|"[CompositeAlternatives]"| q225
    q225 --> q226
    q226 -->|"tok(")")"| q227
    q227 --> q229
    q228 --> q218
    q228 --> q220
    q228 --> q222
    q229 --> q236
    q230 -->|"tok("*")"| q231
    q231 --> q237
    q232 -->|"tok("+")"| q233
    q233 --> q237
    q234 -->|"tok("?")"| q235
    q235 --> q237
    q236 --> q230
    q236 --> q232
    q236 --> q234
    q237 --> q49
```

