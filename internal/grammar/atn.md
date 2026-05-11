# Runtime ATN for grammar

## Grammar

```mermaid
flowchart TD
    q0(["SN:0<br/>RuleStart"])
    q1(["SN:1<br/>RuleStop"])
    q50["SN:50<br/>Basic<br/>"]
    q51["SN:52<br/>Basic<br/>"]
    q52["SN:54<br/>Basic<br/>"]
    q53["SN:56<br/>Basic<br/>"]
    q54["SN:57<br/>Basic<br/>"]
    q55["SN:58<br/>Basic<br/>"]
    q56["SN:59<br/>Basic<br/>"]
    q57["SN:60<br/>Basic<br/>"]
    q58["SN:61<br/>Basic<br/>"]
    q59["SN:62<br/>Basic<br/>"]
    q60["SN:63<br/>Basic<br/>"]
    q61["SN:64<br/>Basic<br/>"]
    q62["SN:65<br/>BlockEnd<br/>"]
    q63{"SN:66<br/>StarLoopEntry<br/><br/>dec=0"}
    q64["SN:67<br/>LoopEnd<br/>"]
    q65["SN:68<br/>StarLoopBack<br/>"]

    q0 --> q50
    q50 -->|"tok(&quot;grammar&quot;)"| q51
    q51 -->|"tok(ID)"| q52
    q52 -->|"tok(&quot;;&quot;)"| q63
    q53 -.->|"[ParserRule]"| q54
    q54 --> q62
    q55 -.->|"[Token]"| q56
    q56 --> q62
    q57 -.->|"[Interface]"| q58
    q58 --> q62
    q59 -.->|"[CompositeRule]"| q60
    q60 --> q62
    q61 --> q53
    q61 --> q55
    q61 --> q57
    q61 --> q59
    q62 --> q65
    q63 --> q61
    q63 --> q64
    q64 --> q1
    q65 --> q63
```

## Interface

```mermaid
flowchart TD
    q2(["SN:2<br/>RuleStart"])
    q3(["SN:3<br/>RuleStop"])
    q66["SN:66<br/>Basic<br/>"]
    q67["SN:68<br/>Basic<br/>"]
    q68["SN:70<br/>Basic<br/>"]
    q69["SN:72<br/>Basic<br/>"]
    q70["SN:74<br/>Basic<br/>"]
    q71["SN:76<br/>Basic<br/>"]
    q72["SN:77<br/>Basic<br/>"]
    q73{"SN:77<br/>StarLoopEntry<br/><br/>dec=1"}
    q74["SN:78<br/>LoopEnd<br/>"]
    q75["SN:79<br/>StarLoopBack<br/>"]
    q76["SN:78<br/>Basic<br/>"]
    q77["SN:80<br/>Basic<br/>"]
    q78["SN:81<br/>Basic<br/>"]
    q79{"SN:82<br/>StarLoopEntry<br/><br/>dec=2"}
    q80["SN:83<br/>LoopEnd<br/>"]
    q81["SN:84<br/>StarLoopBack<br/>"]
    q82["SN:85<br/>Basic<br/>"]
    q83["SN:86<br/>Basic<br/>"]

    q2 --> q66
    q66 -->|"tok(&quot;interface&quot;)"| q67
    q67 -->|"tok(ID)"| q68
    q68 -->|"tok(&quot;extends&quot;)"| q69
    q68 --> q74
    q69 -->|"tok(ID)"| q73
    q70 -->|"tok(&quot;,&quot;)"| q71
    q71 -->|"tok(ID)"| q72
    q72 --> q75
    q73 --> q70
    q73 --> q74
    q74 --> q76
    q75 --> q73
    q76 -->|"tok(&quot;{&quot;)"| q79
    q77 -.->|"[Field]"| q78
    q78 --> q81
    q79 --> q77
    q79 --> q80
    q80 --> q82
    q81 --> q79
    q82 -->|"tok(&quot;}&quot;)"| q83
    q83 --> q3
```

## Field

```mermaid
flowchart TD
    q4(["SN:4<br/>RuleStart"])
    q5(["SN:5<br/>RuleStop"])
    q84["SN:84<br/>Basic<br/>"]
    q85["SN:86<br/>Basic<br/>"]
    q86["SN:87<br/>Basic<br/>"]

    q4 --> q84
    q84 -->|"tok(ID)"| q85
    q85 -.->|"[FieldType]"| q86
    q86 --> q5
```

## FieldType

```mermaid
flowchart TD
    q6(["SN:6<br/>RuleStart"])
    q7(["SN:7<br/>RuleStop"])
    q87["SN:87<br/>Basic<br/>"]
    q88["SN:88<br/>Basic<br/>"]
    q89["SN:89<br/>Basic<br/>"]
    q90["SN:90<br/>Basic<br/>"]
    q91["SN:91<br/>Basic<br/>"]
    q92["SN:92<br/>Basic<br/>"]
    q93["SN:93<br/>Basic<br/>"]
    q94["SN:94<br/>Basic<br/>"]
    q95["SN:95<br/>Basic<br/>"]
    q96["SN:96<br/>BlockEnd<br/>"]

    q6 --> q95
    q87 -.->|"[SimpleType]"| q88
    q88 --> q96
    q89 -.->|"[ReferenceType]"| q90
    q90 --> q96
    q91 -.->|"[ArrayType]"| q92
    q92 --> q96
    q93 -.->|"[PrimitiveType]"| q94
    q94 --> q96
    q95 --> q87
    q95 --> q89
    q95 --> q91
    q95 --> q93
    q96 --> q7
```

## ArrayType

```mermaid
flowchart TD
    q8(["SN:8<br/>RuleStart"])
    q9(["SN:9<br/>RuleStop"])
    q97["SN:97<br/>Basic<br/>"]
    q98["SN:99<br/>Basic<br/>"]
    q99["SN:101<br/>Basic<br/>"]
    q100["SN:102<br/>Basic<br/>"]

    q8 --> q97
    q97 -->|"tok(&quot;[&quot;)"| q98
    q98 -->|"tok(&quot;]&quot;)"| q99
    q99 -.->|"[FieldType]"| q100
    q100 --> q9
```

## ReferenceType

```mermaid
flowchart TD
    q10(["SN:10<br/>RuleStart"])
    q11(["SN:11<br/>RuleStop"])
    q101["SN:101<br/>Basic<br/>"]
    q102["SN:103<br/>Basic<br/>"]
    q103["SN:104<br/>Basic<br/>"]

    q10 --> q101
    q101 -->|"tok(&quot;*&quot;)"| q102
    q102 -->|"tok(ID)"| q103
    q103 --> q11
```

## SimpleType

```mermaid
flowchart TD
    q12(["SN:12<br/>RuleStart"])
    q13(["SN:13<br/>RuleStop"])
    q104["SN:104<br/>Basic<br/>"]
    q105["SN:105<br/>Basic<br/>"]

    q12 --> q104
    q104 -->|"tok(ID)"| q105
    q105 --> q13
```

## PrimitiveType

```mermaid
flowchart TD
    q14(["SN:14<br/>RuleStart"])
    q15(["SN:15<br/>RuleStop"])
    q106["SN:106<br/>Basic<br/>"]
    q107["SN:107<br/>Basic<br/>"]
    q108["SN:108<br/>Basic<br/>"]
    q109["SN:109<br/>Basic<br/>"]
    q110["SN:110<br/>Basic<br/>"]
    q111["SN:111<br/>Basic<br/>"]
    q112["SN:112<br/>Basic<br/>"]
    q113["SN:113<br/>BlockEnd<br/>"]

    q14 --> q112
    q106 -->|"tok(&quot;string&quot;)"| q107
    q107 --> q113
    q108 -->|"tok(&quot;bool&quot;)"| q109
    q109 --> q113
    q110 -->|"tok(&quot;composite&quot;)"| q111
    q111 --> q113
    q112 --> q106
    q112 --> q108
    q112 --> q110
    q113 --> q15
```

## ParserRule

```mermaid
flowchart TD
    q16(["SN:16<br/>RuleStart"])
    q17(["SN:17<br/>RuleStop"])
    q114["SN:114<br/>Basic<br/>"]
    q115["SN:116<br/>Basic<br/>"]
    q116["SN:118<br/>Basic<br/>"]
    q117["SN:119<br/>Basic<br/>"]
    q118["SN:119<br/>Basic<br/>"]
    q119["SN:121<br/>Basic<br/>"]
    q120["SN:123<br/>Basic<br/>"]
    q121["SN:124<br/>Basic<br/>"]

    q16 --> q114
    q114 -->|"tok(ID)"| q115
    q115 -->|"tok(&quot;returns&quot;)"| q116
    q115 --> q117
    q116 -->|"tok(ID)"| q117
    q117 --> q118
    q118 -->|"tok(&quot;:&quot;)"| q119
    q119 -.->|"[Alternatives]"| q120
    q120 -->|"tok(&quot;;&quot;)"| q121
    q121 --> q17
```

## Token

```mermaid
flowchart TD
    q18(["SN:18<br/>RuleStart"])
    q19(["SN:19<br/>RuleStop"])
    q122["SN:122<br/>Basic<br/>"]
    q123["SN:123<br/>Basic<br/>"]
    q124["SN:124<br/>Basic<br/>"]
    q125["SN:125<br/>Basic<br/>"]
    q126["SN:126<br/>Basic<br/>"]
    q127["SN:127<br/>BlockEnd<br/>"]
    q128["SN:128<br/>Basic<br/>"]
    q129["SN:130<br/>Basic<br/>"]
    q130["SN:132<br/>Basic<br/>"]
    q131["SN:134<br/>Basic<br/>"]
    q132["SN:136<br/>Basic<br/>"]
    q133["SN:137<br/>Basic<br/>"]

    q18 --> q126
    q122 -->|"tok(&quot;hidden&quot;)"| q123
    q123 --> q127
    q124 -->|"tok(&quot;comment&quot;)"| q125
    q125 --> q127
    q126 --> q122
    q126 --> q124
    q126 --> q127
    q127 --> q128
    q128 -->|"tok(&quot;token&quot;)"| q129
    q129 -->|"tok(ID)"| q130
    q130 -->|"tok(&quot;:&quot;)"| q131
    q131 -->|"tok(RegexLiteral)"| q132
    q132 -->|"tok(&quot;;&quot;)"| q133
    q133 --> q19
```

## Alternatives

```mermaid
flowchart TD
    q20(["SN:20<br/>RuleStart"])
    q21(["SN:21<br/>RuleStop"])
    q134["SN:134<br/>Basic<br/>"]
    q135["SN:136<br/>Basic<br/>"]
    q136["SN:138<br/>Basic<br/>"]
    q137["SN:139<br/>Basic<br/>"]
    q138{"SN:139<br/>PlusLoopBack<br/><br/>dec=3"}
    q139["SN:140<br/>LoopEnd<br/>"]

    q20 --> q134
    q134 -.->|"[Group]"| q135
    q135 -->|"tok(&quot;|&quot;)"| q136
    q135 --> q139
    q136 -.->|"[Group]"| q137
    q137 --> q138
    q138 --> q135
    q138 --> q139
    q139 --> q21
```

## Group

```mermaid
flowchart TD
    q22(["SN:22<br/>RuleStart"])
    q23(["SN:23<br/>RuleStop"])
    q140["SN:140<br/>Basic<br/>"]
    q141["SN:142<br/>Basic<br/>"]
    q142["SN:143<br/>Basic<br/>"]
    q143{"SN:144<br/>PlusLoopBack<br/><br/>dec=4"}
    q144["SN:145<br/>LoopEnd<br/>"]

    q22 --> q140
    q140 -.->|"[Element]"| q141
    q141 -.->|"[Element]"| q142
    q141 --> q144
    q142 --> q143
    q143 --> q141
    q143 --> q144
    q144 --> q23
```

## Element

```mermaid
flowchart TD
    q24(["SN:24<br/>RuleStart"])
    q25(["SN:25<br/>RuleStop"])
    q145["SN:145<br/>Basic<br/>"]
    q146["SN:146<br/>Basic<br/>"]
    q147["SN:147<br/>Basic<br/>"]
    q148["SN:148<br/>Basic<br/>"]
    q149["SN:149<br/>Basic<br/>"]
    q150["SN:150<br/>Basic<br/>"]
    q151["SN:151<br/>Basic<br/>"]
    q152["SN:152<br/>Basic<br/>"]
    q153["SN:153<br/>Basic<br/>"]
    q154["SN:155<br/>Basic<br/>"]
    q155["SN:157<br/>Basic<br/>"]
    q156["SN:158<br/>Basic<br/>"]
    q157["SN:157<br/>Basic<br/>"]
    q158["SN:158<br/>BlockEnd<br/>"]
    q159["SN:159<br/>Basic<br/>"]
    q160["SN:160<br/>Basic<br/>"]
    q161["SN:161<br/>Basic<br/>"]
    q162["SN:162<br/>Basic<br/>"]
    q163["SN:163<br/>Basic<br/>"]
    q164["SN:164<br/>Basic<br/>"]
    q165["SN:165<br/>Basic<br/>"]
    q166["SN:166<br/>BlockEnd<br/>"]

    q24 --> q157
    q145 -.->|"[Keyword]"| q146
    q146 --> q158
    q147 -.->|"[Assignment]"| q148
    q148 --> q158
    q149 -.->|"[RuleCall]"| q150
    q150 --> q158
    q151 -.->|"[Action]"| q152
    q152 --> q158
    q153 -->|"tok(&quot;(&quot;)"| q154
    q154 -.->|"[Alternatives]"| q155
    q155 -->|"tok(&quot;)&quot;)"| q156
    q156 --> q158
    q157 --> q145
    q157 --> q147
    q157 --> q149
    q157 --> q151
    q157 --> q153
    q158 --> q165
    q159 -->|"tok(&quot;*&quot;)"| q160
    q160 --> q166
    q161 -->|"tok(&quot;+&quot;)"| q162
    q162 --> q166
    q163 -->|"tok(&quot;?&quot;)"| q164
    q164 --> q166
    q165 --> q159
    q165 --> q161
    q165 --> q163
    q165 --> q166
    q166 --> q25
```

## Keyword

```mermaid
flowchart TD
    q26(["SN:26<br/>RuleStart"])
    q27(["SN:27<br/>RuleStop"])
    q167["SN:167<br/>Basic<br/>"]
    q168["SN:168<br/>Basic<br/>"]

    q26 --> q167
    q167 -->|"tok(StringLiteral)"| q168
    q168 --> q27
```

## Assignment

```mermaid
flowchart TD
    q28(["SN:28<br/>RuleStart"])
    q29(["SN:29<br/>RuleStop"])
    q169["SN:169<br/>Basic<br/>"]
    q170["SN:171<br/>Basic<br/>"]
    q171["SN:172<br/>Basic<br/>"]
    q172["SN:173<br/>Basic<br/>"]
    q173["SN:174<br/>Basic<br/>"]
    q174["SN:175<br/>Basic<br/>"]
    q175["SN:176<br/>Basic<br/>"]
    q176["SN:177<br/>Basic<br/>"]
    q177["SN:178<br/>BlockEnd<br/>"]
    q178["SN:179<br/>Basic<br/>"]
    q179["SN:180<br/>Basic<br/>"]

    q28 --> q169
    q169 -->|"tok(ID)"| q176
    q170 -->|"tok(&quot;+=&quot;)"| q171
    q171 --> q177
    q172 -->|"tok(&quot;=&quot;)"| q173
    q173 --> q177
    q174 -->|"tok(&quot;?=&quot;)"| q175
    q175 --> q177
    q176 --> q170
    q176 --> q172
    q176 --> q174
    q177 --> q178
    q178 -.->|"[Assignable]"| q179
    q179 --> q29
```

## Assignable

```mermaid
flowchart TD
    q30(["SN:30<br/>RuleStart"])
    q31(["SN:31<br/>RuleStop"])
    q180["SN:180<br/>Basic<br/>"]
    q181["SN:181<br/>Basic<br/>"]
    q182["SN:182<br/>Basic<br/>"]
    q183["SN:183<br/>Basic<br/>"]
    q184["SN:184<br/>Basic<br/>"]
    q185["SN:185<br/>Basic<br/>"]
    q186["SN:186<br/>Basic<br/>"]
    q187["SN:188<br/>Basic<br/>"]
    q188["SN:190<br/>Basic<br/>"]
    q189["SN:191<br/>Basic<br/>"]
    q190["SN:190<br/>Basic<br/>"]
    q191["SN:191<br/>BlockEnd<br/>"]

    q30 --> q190
    q180 -.->|"[Keyword]"| q181
    q181 --> q191
    q182 -.->|"[RuleCall]"| q183
    q183 --> q191
    q184 -.->|"[CrossRef]"| q185
    q185 --> q191
    q186 -->|"tok(&quot;(&quot;)"| q187
    q187 -.->|"[AssignableAlternatives]"| q188
    q188 -->|"tok(&quot;)&quot;)"| q189
    q189 --> q191
    q190 --> q180
    q190 --> q182
    q190 --> q184
    q190 --> q186
    q191 --> q31
```

## AssignableWithoutAlts

```mermaid
flowchart TD
    q32(["SN:32<br/>RuleStart"])
    q33(["SN:33<br/>RuleStop"])
    q192["SN:192<br/>Basic<br/>"]
    q193["SN:193<br/>Basic<br/>"]
    q194["SN:194<br/>Basic<br/>"]
    q195["SN:195<br/>Basic<br/>"]
    q196["SN:196<br/>Basic<br/>"]
    q197["SN:197<br/>Basic<br/>"]
    q198["SN:198<br/>Basic<br/>"]
    q199["SN:199<br/>BlockEnd<br/>"]

    q32 --> q198
    q192 -.->|"[Keyword]"| q193
    q193 --> q199
    q194 -.->|"[RuleCall]"| q195
    q195 --> q199
    q196 -.->|"[CrossRef]"| q197
    q197 --> q199
    q198 --> q192
    q198 --> q194
    q198 --> q196
    q199 --> q33
```

## AssignableAlternatives

```mermaid
flowchart TD
    q34(["SN:34<br/>RuleStart"])
    q35(["SN:35<br/>RuleStop"])
    q200["SN:200<br/>Basic<br/>"]
    q201["SN:202<br/>Basic<br/>"]
    q202["SN:204<br/>Basic<br/>"]
    q203["SN:205<br/>Basic<br/>"]
    q204{"SN:205<br/>PlusLoopBack<br/><br/>dec=5"}
    q205["SN:206<br/>LoopEnd<br/>"]

    q34 --> q200
    q200 -.->|"[AssignableWithoutAlts]"| q201
    q201 -->|"tok(&quot;|&quot;)"| q202
    q201 --> q205
    q202 -.->|"[AssignableWithoutAlts]"| q203
    q203 --> q204
    q204 --> q201
    q204 --> q205
    q205 --> q35
```

## CrossRef

```mermaid
flowchart TD
    q36(["SN:36<br/>RuleStart"])
    q37(["SN:37<br/>RuleStop"])
    q206["SN:206<br/>Basic<br/>"]
    q207["SN:208<br/>Basic<br/>"]
    q208["SN:210<br/>Basic<br/>"]
    q209["SN:212<br/>Basic<br/>"]
    q210["SN:213<br/>Basic<br/>"]
    q211["SN:213<br/>Basic<br/>"]
    q212["SN:214<br/>Basic<br/>"]

    q36 --> q206
    q206 -->|"tok(&quot;[&quot;)"| q207
    q207 -->|"tok(ID)"| q208
    q208 -->|"tok(&quot;:&quot;)"| q209
    q208 --> q210
    q209 -.->|"[RuleCall]"| q210
    q210 --> q211
    q211 -->|"tok(&quot;]&quot;)"| q212
    q212 --> q37
```

## RuleCall

```mermaid
flowchart TD
    q38(["SN:38<br/>RuleStart"])
    q39(["SN:39<br/>RuleStop"])
    q213["SN:213<br/>Basic<br/>"]
    q214["SN:214<br/>Basic<br/>"]

    q38 --> q213
    q213 -->|"tok(ID)"| q214
    q214 --> q39
```

## Action

```mermaid
flowchart TD
    q40(["SN:40<br/>RuleStart"])
    q41(["SN:41<br/>RuleStop"])
    q215["SN:215<br/>Basic<br/>"]
    q216["SN:217<br/>Basic<br/>"]
    q217["SN:219<br/>Basic<br/>"]
    q218["SN:221<br/>Basic<br/>"]
    q219["SN:223<br/>Basic<br/>"]
    q220["SN:224<br/>Basic<br/>"]
    q221["SN:225<br/>Basic<br/>"]
    q222["SN:226<br/>Basic<br/>"]
    q223["SN:227<br/>Basic<br/>"]
    q224["SN:228<br/>BlockEnd<br/>"]
    q225["SN:229<br/>Basic<br/>"]
    q226["SN:230<br/>Basic<br/>"]
    q227["SN:229<br/>Basic<br/>"]
    q228["SN:230<br/>Basic<br/>"]

    q40 --> q215
    q215 -->|"tok(&quot;{&quot;)"| q216
    q216 -->|"tok(ID)"| q217
    q217 -->|"tok(&quot;.&quot;)"| q218
    q217 --> q226
    q218 -->|"tok(ID)"| q223
    q219 -->|"tok(&quot;+=&quot;)"| q220
    q220 --> q224
    q221 -->|"tok(&quot;=&quot;)"| q222
    q222 --> q224
    q223 --> q219
    q223 --> q221
    q224 --> q225
    q225 -->|"tok(&quot;current&quot;)"| q226
    q226 --> q227
    q227 -->|"tok(&quot;}&quot;)"| q228
    q228 --> q41
```

## CompositeRule

```mermaid
flowchart TD
    q42(["SN:42<br/>RuleStart"])
    q43(["SN:43<br/>RuleStop"])
    q229["SN:229<br/>Basic<br/>"]
    q230["SN:231<br/>Basic<br/>"]
    q231["SN:233<br/>Basic<br/>"]
    q232["SN:235<br/>Basic<br/>"]
    q233["SN:237<br/>Basic<br/>"]
    q234["SN:238<br/>Basic<br/>"]

    q42 --> q229
    q229 -->|"tok(&quot;composite&quot;)"| q230
    q230 -->|"tok(ID)"| q231
    q231 -->|"tok(&quot;:&quot;)"| q232
    q232 -.->|"[CompositeAlternatives]"| q233
    q233 -->|"tok(&quot;;&quot;)"| q234
    q234 --> q43
```

## CompositeAlternatives

```mermaid
flowchart TD
    q44(["SN:44<br/>RuleStart"])
    q45(["SN:45<br/>RuleStop"])
    q235["SN:235<br/>Basic<br/>"]
    q236["SN:237<br/>Basic<br/>"]
    q237["SN:239<br/>Basic<br/>"]
    q238["SN:240<br/>Basic<br/>"]
    q239{"SN:240<br/>PlusLoopBack<br/><br/>dec=6"}
    q240["SN:241<br/>LoopEnd<br/>"]

    q44 --> q235
    q235 -.->|"[CompositeGroup]"| q236
    q236 -->|"tok(&quot;|&quot;)"| q237
    q236 --> q240
    q237 -.->|"[CompositeGroup]"| q238
    q238 --> q239
    q239 --> q236
    q239 --> q240
    q240 --> q45
```

## CompositeGroup

```mermaid
flowchart TD
    q46(["SN:46<br/>RuleStart"])
    q47(["SN:47<br/>RuleStop"])
    q241["SN:241<br/>Basic<br/>"]
    q242["SN:243<br/>Basic<br/>"]
    q243["SN:244<br/>Basic<br/>"]
    q244{"SN:245<br/>PlusLoopBack<br/><br/>dec=7"}
    q245["SN:246<br/>LoopEnd<br/>"]

    q46 --> q241
    q241 -.->|"[CompositeElement]"| q242
    q242 -.->|"[CompositeElement]"| q243
    q242 --> q245
    q243 --> q244
    q244 --> q242
    q244 --> q245
    q245 --> q47
```

## CompositeElement

```mermaid
flowchart TD
    q48(["SN:48<br/>RuleStart"])
    q49(["SN:49<br/>RuleStop"])
    q246["SN:246<br/>Basic<br/>"]
    q247["SN:247<br/>Basic<br/>"]
    q248["SN:248<br/>Basic<br/>"]
    q249["SN:249<br/>Basic<br/>"]
    q250["SN:250<br/>Basic<br/>"]
    q251["SN:252<br/>Basic<br/>"]
    q252["SN:254<br/>Basic<br/>"]
    q253["SN:255<br/>Basic<br/>"]
    q254["SN:254<br/>Basic<br/>"]
    q255["SN:255<br/>BlockEnd<br/>"]
    q256["SN:256<br/>Basic<br/>"]
    q257["SN:257<br/>Basic<br/>"]
    q258["SN:258<br/>Basic<br/>"]
    q259["SN:259<br/>Basic<br/>"]
    q260["SN:260<br/>Basic<br/>"]
    q261["SN:261<br/>Basic<br/>"]
    q262["SN:262<br/>Basic<br/>"]
    q263["SN:263<br/>BlockEnd<br/>"]

    q48 --> q254
    q246 -.->|"[Keyword]"| q247
    q247 --> q255
    q248 -.->|"[RuleCall]"| q249
    q249 --> q255
    q250 -->|"tok(&quot;(&quot;)"| q251
    q251 -.->|"[CompositeAlternatives]"| q252
    q252 -->|"tok(&quot;)&quot;)"| q253
    q253 --> q255
    q254 --> q246
    q254 --> q248
    q254 --> q250
    q255 --> q262
    q256 -->|"tok(&quot;*&quot;)"| q257
    q257 --> q263
    q258 -->|"tok(&quot;+&quot;)"| q259
    q259 --> q263
    q260 -->|"tok(&quot;?&quot;)"| q261
    q261 --> q263
    q262 --> q256
    q262 --> q258
    q262 --> q260
    q262 --> q263
    q263 --> q49
```

