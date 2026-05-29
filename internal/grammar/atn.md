# Runtime ATN for grammar

## Grammar

```mermaid
flowchart TD
    q0(["SN:0<br/>RuleStart"])
    q1(["SN:1<br/>RuleStop"])
    q52["SN:52<br/>Basic<br/>"]
    q53["SN:54<br/>Basic<br/>"]
    q54["SN:56<br/>Basic<br/>"]
    q55["SN:58<br/>Basic<br/>"]
    q56["SN:59<br/>Basic<br/>"]
    q57["SN:60<br/>Basic<br/>"]
    q58["SN:61<br/>Basic<br/>"]
    q59["SN:62<br/>Basic<br/>"]
    q60["SN:63<br/>Basic<br/>"]
    q61["SN:64<br/>Basic<br/>"]
    q62["SN:65<br/>Basic<br/>"]
    q63["SN:66<br/>Basic<br/>"]
    q64["SN:67<br/>Basic<br/>"]
    q65["SN:68<br/>Basic<br/>"]
    q66["SN:69<br/>BlockEnd<br/>"]
    q67{"SN:70<br/>StarLoopEntry<br/><br/>dec=0"}
    q68["SN:71<br/>LoopEnd<br/>"]
    q69["SN:72<br/>StarLoopBack<br/>"]

    q0 --> q52
    q52 -->|"tok(&quot;grammar&quot;)"| q53
    q53 -->|"tok(ID)"| q54
    q54 -->|"tok(&quot;;&quot;)"| q67
    q55 -.->|"[ParserRule]"| q56
    q56 --> q66
    q57 -.->|"[Token]"| q58
    q58 --> q66
    q59 -.->|"[TokenGroup]"| q60
    q60 --> q66
    q61 -.->|"[Interface]"| q62
    q62 --> q66
    q63 -.->|"[CompositeRule]"| q64
    q64 --> q66
    q65 --> q55
    q65 --> q57
    q65 --> q59
    q65 --> q61
    q65 --> q63
    q66 --> q69
    q67 --> q65
    q67 --> q68
    q68 --> q1
    q69 --> q67
```

## Interface

```mermaid
flowchart TD
    q2(["SN:2<br/>RuleStart"])
    q3(["SN:3<br/>RuleStop"])
    q70["SN:70<br/>Basic<br/>"]
    q71["SN:72<br/>Basic<br/>"]
    q72["SN:74<br/>Basic<br/>"]
    q73["SN:76<br/>Basic<br/>"]
    q74["SN:78<br/>Basic<br/>"]
    q75["SN:80<br/>Basic<br/>"]
    q76["SN:81<br/>Basic<br/>"]
    q77{"SN:81<br/>StarLoopEntry<br/><br/>dec=1"}
    q78["SN:82<br/>LoopEnd<br/>"]
    q79["SN:83<br/>StarLoopBack<br/>"]
    q80["SN:82<br/>Basic<br/>"]
    q81["SN:84<br/>Basic<br/>"]
    q82["SN:85<br/>Basic<br/>"]
    q83{"SN:86<br/>StarLoopEntry<br/><br/>dec=2"}
    q84["SN:87<br/>LoopEnd<br/>"]
    q85["SN:88<br/>StarLoopBack<br/>"]
    q86["SN:89<br/>Basic<br/>"]
    q87["SN:90<br/>Basic<br/>"]

    q2 --> q70
    q70 -->|"tok(&quot;interface&quot;)"| q71
    q71 -->|"tok(ID)"| q72
    q72 -->|"tok(&quot;extends&quot;)"| q73
    q72 --> q78
    q73 -->|"tok(ID)"| q77
    q74 -->|"tok(&quot;,&quot;)"| q75
    q75 -->|"tok(ID)"| q76
    q76 --> q79
    q77 --> q74
    q77 --> q78
    q78 --> q80
    q79 --> q77
    q80 -->|"tok(&quot;{&quot;)"| q83
    q81 -.->|"[Field]"| q82
    q82 --> q85
    q83 --> q81
    q83 --> q84
    q84 --> q86
    q85 --> q83
    q86 -->|"tok(&quot;}&quot;)"| q87
    q87 --> q3
```

## Field

```mermaid
flowchart TD
    q4(["SN:4<br/>RuleStart"])
    q5(["SN:5<br/>RuleStop"])
    q88["SN:88<br/>Basic<br/>"]
    q89["SN:90<br/>Basic<br/>"]
    q90["SN:91<br/>Basic<br/>"]

    q4 --> q88
    q88 -->|"tok(ID)"| q89
    q89 -.->|"[FieldType]"| q90
    q90 --> q5
```

## FieldType

```mermaid
flowchart TD
    q6(["SN:6<br/>RuleStart"])
    q7(["SN:7<br/>RuleStop"])
    q91["SN:91<br/>Basic<br/>"]
    q92["SN:92<br/>Basic<br/>"]
    q93["SN:93<br/>Basic<br/>"]
    q94["SN:94<br/>Basic<br/>"]
    q95["SN:95<br/>Basic<br/>"]
    q96["SN:96<br/>Basic<br/>"]
    q97["SN:97<br/>Basic<br/>"]
    q98["SN:98<br/>Basic<br/>"]
    q99["SN:99<br/>Basic<br/>"]
    q100["SN:100<br/>BlockEnd<br/>"]

    q6 --> q99
    q91 -.->|"[SimpleType]"| q92
    q92 --> q100
    q93 -.->|"[ReferenceType]"| q94
    q94 --> q100
    q95 -.->|"[ArrayType]"| q96
    q96 --> q100
    q97 -.->|"[PrimitiveType]"| q98
    q98 --> q100
    q99 --> q91
    q99 --> q93
    q99 --> q95
    q99 --> q97
    q100 --> q7
```

## ArrayType

```mermaid
flowchart TD
    q8(["SN:8<br/>RuleStart"])
    q9(["SN:9<br/>RuleStop"])
    q101["SN:101<br/>Basic<br/>"]
    q102["SN:103<br/>Basic<br/>"]
    q103["SN:105<br/>Basic<br/>"]
    q104["SN:106<br/>Basic<br/>"]

    q8 --> q101
    q101 -->|"tok(&quot;[&quot;)"| q102
    q102 -->|"tok(&quot;]&quot;)"| q103
    q103 -.->|"[FieldType]"| q104
    q104 --> q9
```

## ReferenceType

```mermaid
flowchart TD
    q10(["SN:10<br/>RuleStart"])
    q11(["SN:11<br/>RuleStop"])
    q105["SN:105<br/>Basic<br/>"]
    q106["SN:107<br/>Basic<br/>"]
    q107["SN:108<br/>Basic<br/>"]

    q10 --> q105
    q105 -->|"tok(&quot;*&quot;)"| q106
    q106 -->|"tok(ID)"| q107
    q107 --> q11
```

## SimpleType

```mermaid
flowchart TD
    q12(["SN:12<br/>RuleStart"])
    q13(["SN:13<br/>RuleStop"])
    q108["SN:108<br/>Basic<br/>"]
    q109["SN:109<br/>Basic<br/>"]

    q12 --> q108
    q108 -->|"tok(ID)"| q109
    q109 --> q13
```

## PrimitiveType

```mermaid
flowchart TD
    q14(["SN:14<br/>RuleStart"])
    q15(["SN:15<br/>RuleStop"])
    q110["SN:110<br/>Basic<br/>"]
    q111["SN:111<br/>Basic<br/>"]
    q112["SN:112<br/>Basic<br/>"]
    q113["SN:113<br/>Basic<br/>"]
    q114["SN:114<br/>Basic<br/>"]
    q115["SN:115<br/>Basic<br/>"]
    q116["SN:116<br/>Basic<br/>"]
    q117["SN:117<br/>BlockEnd<br/>"]

    q14 --> q116
    q110 -->|"tok(&quot;string&quot;)"| q111
    q111 --> q117
    q112 -->|"tok(&quot;bool&quot;)"| q113
    q113 --> q117
    q114 -->|"tok(&quot;composite&quot;)"| q115
    q115 --> q117
    q116 --> q110
    q116 --> q112
    q116 --> q114
    q117 --> q15
```

## ParserRule

```mermaid
flowchart TD
    q16(["SN:16<br/>RuleStart"])
    q17(["SN:17<br/>RuleStop"])
    q118["SN:118<br/>Basic<br/>"]
    q119["SN:120<br/>Basic<br/>"]
    q120["SN:122<br/>Basic<br/>"]
    q121["SN:123<br/>Basic<br/>"]
    q122["SN:123<br/>Basic<br/>"]
    q123["SN:125<br/>Basic<br/>"]
    q124["SN:127<br/>Basic<br/>"]
    q125["SN:128<br/>Basic<br/>"]

    q16 --> q118
    q118 -->|"tok(ID)"| q119
    q119 -->|"tok(&quot;returns&quot;)"| q120
    q119 --> q121
    q120 -->|"tok(ID)"| q121
    q121 --> q122
    q122 -->|"tok(&quot;:&quot;)"| q123
    q123 -.->|"[Alternatives]"| q124
    q124 -->|"tok(&quot;;&quot;)"| q125
    q125 --> q17
```

## Token

```mermaid
flowchart TD
    q18(["SN:18<br/>RuleStart"])
    q19(["SN:19<br/>RuleStop"])
    q126["SN:126<br/>Basic<br/>"]
    q127["SN:127<br/>Basic<br/>"]
    q128["SN:128<br/>Basic<br/>"]
    q129["SN:129<br/>Basic<br/>"]
    q130["SN:130<br/>Basic<br/>"]
    q131["SN:131<br/>BlockEnd<br/>"]
    q132["SN:132<br/>Basic<br/>"]
    q133["SN:134<br/>Basic<br/>"]
    q134["SN:136<br/>Basic<br/>"]
    q135["SN:138<br/>Basic<br/>"]
    q136["SN:140<br/>Basic<br/>"]
    q137["SN:141<br/>Basic<br/>"]

    q18 --> q130
    q126 -->|"tok(&quot;hidden&quot;)"| q127
    q127 --> q131
    q128 -->|"tok(&quot;comment&quot;)"| q129
    q129 --> q131
    q130 --> q126
    q130 --> q128
    q130 --> q131
    q131 --> q132
    q132 -->|"tok(&quot;token&quot;)"| q133
    q133 -->|"tok(ID)"| q134
    q134 -->|"tok(&quot;:&quot;)"| q135
    q135 -->|"tok(RegexLiteral)"| q136
    q136 -->|"tok(&quot;;&quot;)"| q137
    q137 --> q19
```

## TokenGroup

```mermaid
flowchart TD
    q20(["SN:20<br/>RuleStart"])
    q21(["SN:21<br/>RuleStop"])
    q138["SN:138<br/>Basic<br/>"]
    q139["SN:140<br/>Basic<br/>"]
    q140["SN:142<br/>Basic<br/>"]
    q141["SN:144<br/>Basic<br/>"]
    q142["SN:146<br/>Basic<br/>"]
    q143["SN:147<br/>Basic<br/>"]
    q144["SN:148<br/>Basic<br/>"]
    q145["SN:149<br/>Basic<br/>"]
    q146["SN:150<br/>Basic<br/>"]
    q147["SN:151<br/>Basic<br/>"]
    q148["SN:152<br/>Basic<br/>"]
    q149["SN:153<br/>BlockEnd<br/>"]
    q150{"SN:154<br/>StarLoopEntry<br/><br/>dec=3"}
    q151["SN:155<br/>LoopEnd<br/>"]
    q152["SN:156<br/>StarLoopBack<br/>"]
    q153["SN:157<br/>Basic<br/>"]
    q154["SN:158<br/>Basic<br/>"]

    q20 --> q138
    q138 -->|"tok(&quot;token&quot;)"| q139
    q139 -->|"tok(&quot;group&quot;)"| q140
    q140 -->|"tok(ID)"| q141
    q141 -->|"tok(&quot;{&quot;)"| q150
    q142 -->|"tok(ID)"| q143
    q143 --> q149
    q144 -->|"tok(RegexLiteral)"| q145
    q145 --> q149
    q146 -->|"tok(StringLiteral)"| q147
    q147 --> q149
    q148 --> q142
    q148 --> q144
    q148 --> q146
    q149 --> q152
    q150 --> q148
    q150 --> q151
    q151 --> q153
    q152 --> q150
    q153 -->|"tok(&quot;}&quot;)"| q154
    q154 --> q21
```

## Alternatives

```mermaid
flowchart TD
    q22(["SN:22<br/>RuleStart"])
    q23(["SN:23<br/>RuleStop"])
    q155["SN:155<br/>Basic<br/>"]
    q156["SN:157<br/>Basic<br/>"]
    q157["SN:159<br/>Basic<br/>"]
    q158["SN:160<br/>Basic<br/>"]
    q159{"SN:160<br/>PlusLoopBack<br/><br/>dec=4"}
    q160["SN:161<br/>LoopEnd<br/>"]

    q22 --> q155
    q155 -.->|"[Group]"| q156
    q156 -->|"tok(&quot;|&quot;)"| q157
    q156 --> q160
    q157 -.->|"[Group]"| q158
    q158 --> q159
    q159 --> q156
    q159 --> q160
    q160 --> q23
```

## Group

```mermaid
flowchart TD
    q24(["SN:24<br/>RuleStart"])
    q25(["SN:25<br/>RuleStop"])
    q161["SN:161<br/>Basic<br/>"]
    q162["SN:163<br/>Basic<br/>"]
    q163["SN:164<br/>Basic<br/>"]
    q164{"SN:165<br/>PlusLoopBack<br/><br/>dec=5"}
    q165["SN:166<br/>LoopEnd<br/>"]

    q24 --> q161
    q161 -.->|"[Element]"| q162
    q162 -.->|"[Element]"| q163
    q162 --> q165
    q163 --> q164
    q164 --> q162
    q164 --> q165
    q165 --> q25
```

## Element

```mermaid
flowchart TD
    q26(["SN:26<br/>RuleStart"])
    q27(["SN:27<br/>RuleStop"])
    q166["SN:166<br/>Basic<br/>"]
    q167["SN:167<br/>Basic<br/>"]
    q168["SN:168<br/>Basic<br/>"]
    q169["SN:169<br/>Basic<br/>"]
    q170["SN:170<br/>Basic<br/>"]
    q171["SN:171<br/>Basic<br/>"]
    q172["SN:172<br/>Basic<br/>"]
    q173["SN:173<br/>Basic<br/>"]
    q174["SN:174<br/>Basic<br/>"]
    q175["SN:176<br/>Basic<br/>"]
    q176["SN:178<br/>Basic<br/>"]
    q177["SN:179<br/>Basic<br/>"]
    q178["SN:178<br/>Basic<br/>"]
    q179["SN:179<br/>BlockEnd<br/>"]
    q180["SN:180<br/>Basic<br/>"]
    q181["SN:181<br/>Basic<br/>"]
    q182["SN:182<br/>Basic<br/>"]
    q183["SN:183<br/>Basic<br/>"]
    q184["SN:184<br/>Basic<br/>"]
    q185["SN:185<br/>Basic<br/>"]
    q186["SN:186<br/>Basic<br/>"]
    q187["SN:187<br/>BlockEnd<br/>"]

    q26 --> q178
    q166 -.->|"[Keyword]"| q167
    q167 --> q179
    q168 -.->|"[Assignment]"| q169
    q169 --> q179
    q170 -.->|"[RuleCall]"| q171
    q171 --> q179
    q172 -.->|"[Action]"| q173
    q173 --> q179
    q174 -->|"tok(&quot;(&quot;)"| q175
    q175 -.->|"[Alternatives]"| q176
    q176 -->|"tok(&quot;)&quot;)"| q177
    q177 --> q179
    q178 --> q166
    q178 --> q168
    q178 --> q170
    q178 --> q172
    q178 --> q174
    q179 --> q186
    q180 -->|"tok(&quot;*&quot;)"| q181
    q181 --> q187
    q182 -->|"tok(&quot;+&quot;)"| q183
    q183 --> q187
    q184 -->|"tok(&quot;?&quot;)"| q185
    q185 --> q187
    q186 --> q180
    q186 --> q182
    q186 --> q184
    q186 --> q187
    q187 --> q27
```

## Keyword

```mermaid
flowchart TD
    q28(["SN:28<br/>RuleStart"])
    q29(["SN:29<br/>RuleStop"])
    q188["SN:188<br/>Basic<br/>"]
    q189["SN:189<br/>Basic<br/>"]

    q28 --> q188
    q188 -->|"tok(StringLiteral)"| q189
    q189 --> q29
```

## Assignment

```mermaid
flowchart TD
    q30(["SN:30<br/>RuleStart"])
    q31(["SN:31<br/>RuleStop"])
    q190["SN:190<br/>Basic<br/>"]
    q191["SN:192<br/>Basic<br/>"]
    q192["SN:193<br/>Basic<br/>"]
    q193["SN:194<br/>Basic<br/>"]
    q194["SN:195<br/>Basic<br/>"]
    q195["SN:196<br/>Basic<br/>"]
    q196["SN:197<br/>Basic<br/>"]
    q197["SN:198<br/>Basic<br/>"]
    q198["SN:199<br/>BlockEnd<br/>"]
    q199["SN:200<br/>Basic<br/>"]
    q200["SN:201<br/>Basic<br/>"]

    q30 --> q190
    q190 -->|"tok(ID)"| q197
    q191 -->|"tok(&quot;+=&quot;)"| q192
    q192 --> q198
    q193 -->|"tok(&quot;=&quot;)"| q194
    q194 --> q198
    q195 -->|"tok(&quot;?=&quot;)"| q196
    q196 --> q198
    q197 --> q191
    q197 --> q193
    q197 --> q195
    q198 --> q199
    q199 -.->|"[Assignable]"| q200
    q200 --> q31
```

## Assignable

```mermaid
flowchart TD
    q32(["SN:32<br/>RuleStart"])
    q33(["SN:33<br/>RuleStop"])
    q201["SN:201<br/>Basic<br/>"]
    q202["SN:202<br/>Basic<br/>"]
    q203["SN:203<br/>Basic<br/>"]
    q204["SN:204<br/>Basic<br/>"]
    q205["SN:205<br/>Basic<br/>"]
    q206["SN:206<br/>Basic<br/>"]
    q207["SN:207<br/>Basic<br/>"]
    q208["SN:209<br/>Basic<br/>"]
    q209["SN:211<br/>Basic<br/>"]
    q210["SN:212<br/>Basic<br/>"]
    q211["SN:211<br/>Basic<br/>"]
    q212["SN:212<br/>BlockEnd<br/>"]

    q32 --> q211
    q201 -.->|"[Keyword]"| q202
    q202 --> q212
    q203 -.->|"[RuleCall]"| q204
    q204 --> q212
    q205 -.->|"[CrossRef]"| q206
    q206 --> q212
    q207 -->|"tok(&quot;(&quot;)"| q208
    q208 -.->|"[AssignableAlternatives]"| q209
    q209 -->|"tok(&quot;)&quot;)"| q210
    q210 --> q212
    q211 --> q201
    q211 --> q203
    q211 --> q205
    q211 --> q207
    q212 --> q33
```

## AssignableWithoutAlts

```mermaid
flowchart TD
    q34(["SN:34<br/>RuleStart"])
    q35(["SN:35<br/>RuleStop"])
    q213["SN:213<br/>Basic<br/>"]
    q214["SN:214<br/>Basic<br/>"]
    q215["SN:215<br/>Basic<br/>"]
    q216["SN:216<br/>Basic<br/>"]
    q217["SN:217<br/>Basic<br/>"]
    q218["SN:218<br/>Basic<br/>"]
    q219["SN:219<br/>Basic<br/>"]
    q220["SN:220<br/>BlockEnd<br/>"]

    q34 --> q219
    q213 -.->|"[Keyword]"| q214
    q214 --> q220
    q215 -.->|"[RuleCall]"| q216
    q216 --> q220
    q217 -.->|"[CrossRef]"| q218
    q218 --> q220
    q219 --> q213
    q219 --> q215
    q219 --> q217
    q220 --> q35
```

## AssignableAlternatives

```mermaid
flowchart TD
    q36(["SN:36<br/>RuleStart"])
    q37(["SN:37<br/>RuleStop"])
    q221["SN:221<br/>Basic<br/>"]
    q222["SN:223<br/>Basic<br/>"]
    q223["SN:225<br/>Basic<br/>"]
    q224["SN:226<br/>Basic<br/>"]
    q225{"SN:226<br/>PlusLoopBack<br/><br/>dec=6"}
    q226["SN:227<br/>LoopEnd<br/>"]

    q36 --> q221
    q221 -.->|"[AssignableWithoutAlts]"| q222
    q222 -->|"tok(&quot;|&quot;)"| q223
    q222 --> q226
    q223 -.->|"[AssignableWithoutAlts]"| q224
    q224 --> q225
    q225 --> q222
    q225 --> q226
    q226 --> q37
```

## CrossRef

```mermaid
flowchart TD
    q38(["SN:38<br/>RuleStart"])
    q39(["SN:39<br/>RuleStop"])
    q227["SN:227<br/>Basic<br/>"]
    q228["SN:229<br/>Basic<br/>"]
    q229["SN:231<br/>Basic<br/>"]
    q230["SN:233<br/>Basic<br/>"]
    q231["SN:234<br/>Basic<br/>"]
    q232["SN:234<br/>Basic<br/>"]
    q233["SN:235<br/>Basic<br/>"]

    q38 --> q227
    q227 -->|"tok(&quot;[&quot;)"| q228
    q228 -->|"tok(ID)"| q229
    q229 -->|"tok(&quot;:&quot;)"| q230
    q229 --> q231
    q230 -.->|"[RuleCall]"| q231
    q231 --> q232
    q232 -->|"tok(&quot;]&quot;)"| q233
    q233 --> q39
```

## RuleCall

```mermaid
flowchart TD
    q40(["SN:40<br/>RuleStart"])
    q41(["SN:41<br/>RuleStop"])
    q234["SN:234<br/>Basic<br/>"]
    q235["SN:235<br/>Basic<br/>"]

    q40 --> q234
    q234 -->|"tok(ID)"| q235
    q235 --> q41
```

## Action

```mermaid
flowchart TD
    q42(["SN:42<br/>RuleStart"])
    q43(["SN:43<br/>RuleStop"])
    q236["SN:236<br/>Basic<br/>"]
    q237["SN:238<br/>Basic<br/>"]
    q238["SN:240<br/>Basic<br/>"]
    q239["SN:242<br/>Basic<br/>"]
    q240["SN:244<br/>Basic<br/>"]
    q241["SN:245<br/>Basic<br/>"]
    q242["SN:246<br/>Basic<br/>"]
    q243["SN:247<br/>Basic<br/>"]
    q244["SN:248<br/>Basic<br/>"]
    q245["SN:249<br/>BlockEnd<br/>"]
    q246["SN:250<br/>Basic<br/>"]
    q247["SN:251<br/>Basic<br/>"]
    q248["SN:250<br/>Basic<br/>"]
    q249["SN:251<br/>Basic<br/>"]

    q42 --> q236
    q236 -->|"tok(&quot;{&quot;)"| q237
    q237 -->|"tok(ID)"| q238
    q238 -->|"tok(&quot;.&quot;)"| q239
    q238 --> q247
    q239 -->|"tok(ID)"| q244
    q240 -->|"tok(&quot;+=&quot;)"| q241
    q241 --> q245
    q242 -->|"tok(&quot;=&quot;)"| q243
    q243 --> q245
    q244 --> q240
    q244 --> q242
    q245 --> q246
    q246 -->|"tok(&quot;current&quot;)"| q247
    q247 --> q248
    q248 -->|"tok(&quot;}&quot;)"| q249
    q249 --> q43
```

## CompositeRule

```mermaid
flowchart TD
    q44(["SN:44<br/>RuleStart"])
    q45(["SN:45<br/>RuleStop"])
    q250["SN:250<br/>Basic<br/>"]
    q251["SN:252<br/>Basic<br/>"]
    q252["SN:254<br/>Basic<br/>"]
    q253["SN:256<br/>Basic<br/>"]
    q254["SN:258<br/>Basic<br/>"]
    q255["SN:259<br/>Basic<br/>"]

    q44 --> q250
    q250 -->|"tok(&quot;composite&quot;)"| q251
    q251 -->|"tok(ID)"| q252
    q252 -->|"tok(&quot;:&quot;)"| q253
    q253 -.->|"[CompositeAlternatives]"| q254
    q254 -->|"tok(&quot;;&quot;)"| q255
    q255 --> q45
```

## CompositeAlternatives

```mermaid
flowchart TD
    q46(["SN:46<br/>RuleStart"])
    q47(["SN:47<br/>RuleStop"])
    q256["SN:256<br/>Basic<br/>"]
    q257["SN:258<br/>Basic<br/>"]
    q258["SN:260<br/>Basic<br/>"]
    q259["SN:261<br/>Basic<br/>"]
    q260{"SN:261<br/>PlusLoopBack<br/><br/>dec=7"}
    q261["SN:262<br/>LoopEnd<br/>"]

    q46 --> q256
    q256 -.->|"[CompositeGroup]"| q257
    q257 -->|"tok(&quot;|&quot;)"| q258
    q257 --> q261
    q258 -.->|"[CompositeGroup]"| q259
    q259 --> q260
    q260 --> q257
    q260 --> q261
    q261 --> q47
```

## CompositeGroup

```mermaid
flowchart TD
    q48(["SN:48<br/>RuleStart"])
    q49(["SN:49<br/>RuleStop"])
    q262["SN:262<br/>Basic<br/>"]
    q263["SN:264<br/>Basic<br/>"]
    q264["SN:265<br/>Basic<br/>"]
    q265{"SN:266<br/>PlusLoopBack<br/><br/>dec=8"}
    q266["SN:267<br/>LoopEnd<br/>"]

    q48 --> q262
    q262 -.->|"[CompositeElement]"| q263
    q263 -.->|"[CompositeElement]"| q264
    q263 --> q266
    q264 --> q265
    q265 --> q263
    q265 --> q266
    q266 --> q49
```

## CompositeElement

```mermaid
flowchart TD
    q50(["SN:50<br/>RuleStart"])
    q51(["SN:51<br/>RuleStop"])
    q267["SN:267<br/>Basic<br/>"]
    q268["SN:268<br/>Basic<br/>"]
    q269["SN:269<br/>Basic<br/>"]
    q270["SN:270<br/>Basic<br/>"]
    q271["SN:271<br/>Basic<br/>"]
    q272["SN:273<br/>Basic<br/>"]
    q273["SN:275<br/>Basic<br/>"]
    q274["SN:276<br/>Basic<br/>"]
    q275["SN:275<br/>Basic<br/>"]
    q276["SN:276<br/>BlockEnd<br/>"]
    q277["SN:277<br/>Basic<br/>"]
    q278["SN:278<br/>Basic<br/>"]
    q279["SN:279<br/>Basic<br/>"]
    q280["SN:280<br/>Basic<br/>"]
    q281["SN:281<br/>Basic<br/>"]
    q282["SN:282<br/>Basic<br/>"]
    q283["SN:283<br/>Basic<br/>"]
    q284["SN:284<br/>BlockEnd<br/>"]

    q50 --> q275
    q267 -.->|"[Keyword]"| q268
    q268 --> q276
    q269 -.->|"[RuleCall]"| q270
    q270 --> q276
    q271 -->|"tok(&quot;(&quot;)"| q272
    q272 -.->|"[CompositeAlternatives]"| q273
    q273 -->|"tok(&quot;)&quot;)"| q274
    q274 --> q276
    q275 --> q267
    q275 --> q269
    q275 --> q271
    q276 --> q283
    q277 -->|"tok(&quot;*&quot;)"| q278
    q278 --> q284
    q279 -->|"tok(&quot;+&quot;)"| q280
    q280 --> q284
    q281 -->|"tok(&quot;?&quot;)"| q282
    q282 --> q284
    q283 --> q277
    q283 --> q279
    q283 --> q281
    q283 --> q284
    q284 --> q51
```

