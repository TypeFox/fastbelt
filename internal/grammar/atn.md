# Runtime ATN for grammar

## Grammar

```mermaid
flowchart TD
    q0(["SN:0<br/>RuleStart"])
    q1(["SN:1<br/>RuleStop"])
    q50["SN:50<br/>Basic<br/>Terminal #1"]
    q51["SN:52<br/>Basic<br/>Terminal #2"]
    q52["SN:54<br/>Basic<br/>Terminal #3"]
    q53{"SN:56<br/>StarBlockStart<br/>Repetition #1<br/>dec=0"}
    q54{"SN:57<br/>Basic<br/>Alternation #1<br/>dec=1"}
    q55["SN:58<br/>Basic<br/>NonTerminal #1"]
    q56["SN:59<br/>Basic<br/>NonTerminal #1"]
    q57["SN:60<br/>Basic<br/>NonTerminal #2"]
    q58["SN:61<br/>Basic<br/>NonTerminal #2"]
    q59["SN:62<br/>Basic<br/>NonTerminal #3"]
    q60["SN:63<br/>Basic<br/>NonTerminal #3"]
    q61["SN:64<br/>Basic<br/>NonTerminal #4"]
    q62["SN:65<br/>Basic<br/>NonTerminal #4"]
    q63["SN:66<br/>BlockEnd<br/>Alternation #1"]
    q64["SN:67<br/>BlockEnd<br/>Repetition #1"]
    q65{"SN:68<br/>StarLoopEntry<br/>Repetition #1<br/>dec=2"}
    q66["SN:69<br/>LoopEnd<br/>Repetition #1"]
    q67["SN:70<br/>StarLoopBack<br/>Repetition #1"]

    q0 --> q50
    q50 -->|"tok('grammar')"| q51
    q51 -->|"tok(ID)"| q52
    q52 -->|"tok(';')"| q65
    q53 --> q54
    q54 --> q55
    q54 --> q57
    q54 --> q59
    q54 --> q61
    q55 -.->|"[ParserRule]"| q56
    q56 --> q63
    q57 -.->|"[Token]"| q58
    q58 --> q63
    q59 -.->|"[Interface]"| q60
    q60 --> q63
    q61 -.->|"[CompositeRule]"| q62
    q62 --> q63
    q63 --> q64
    q64 --> q67
    q65 --> q53
    q65 --> q66
    q66 --> q1
    q67 --> q65
```

## Interface

```mermaid
flowchart TD
    q2(["SN:2<br/>RuleStart"])
    q3(["SN:3<br/>RuleStop"])
    q68["SN:68<br/>Basic<br/>Terminal #1"]
    q69["SN:70<br/>Basic<br/>Terminal #2"]
    q70{"SN:72<br/>Basic<br/>Option #1<br/>dec=3"}
    q71["SN:73<br/>Basic<br/>Terminal #3"]
    q72["SN:75<br/>Basic<br/>Terminal #4"]
    q73{"SN:77<br/>StarBlockStart<br/>Repetition #1<br/>dec=4"}
    q74["SN:78<br/>Basic<br/>Terminal #5"]
    q75["SN:80<br/>Basic<br/>Terminal #6"]
    q76["SN:81<br/>Basic<br/>Terminal #6"]
    q77["SN:81<br/>BlockEnd<br/>Repetition #1"]
    q78{"SN:82<br/>StarLoopEntry<br/>Repetition #1<br/>dec=5"}
    q79["SN:83<br/>LoopEnd<br/>Repetition #1"]
    q80["SN:84<br/>StarLoopBack<br/>Repetition #1"]
    q81["SN:83<br/>BlockEnd<br/>Option #1"]
    q82["SN:84<br/>Basic<br/>Terminal #7"]
    q83{"SN:86<br/>StarBlockStart<br/>Repetition #2<br/>dec=6"}
    q84["SN:87<br/>Basic<br/>NonTerminal #1"]
    q85["SN:88<br/>Basic<br/>NonTerminal #1"]
    q86["SN:89<br/>BlockEnd<br/>Repetition #2"]
    q87{"SN:90<br/>StarLoopEntry<br/>Repetition #2<br/>dec=7"}
    q88["SN:91<br/>LoopEnd<br/>Repetition #2"]
    q89["SN:92<br/>StarLoopBack<br/>Repetition #2"]
    q90["SN:93<br/>Basic<br/>Terminal #8"]
    q91["SN:94<br/>Basic<br/>Terminal #8"]

    q2 --> q68
    q68 -->|"tok('interface')"| q69
    q69 -->|"tok(ID)"| q70
    q70 --> q71
    q70 --> q81
    q71 -->|"tok('extends')"| q72
    q72 -->|"tok(ID)"| q78
    q73 --> q74
    q74 -->|"tok(',')"| q75
    q75 -->|"tok(ID)"| q76
    q76 --> q77
    q77 --> q80
    q78 --> q73
    q78 --> q79
    q79 --> q81
    q80 --> q78
    q81 --> q82
    q82 -->|"tok('{')"| q87
    q83 --> q84
    q84 -.->|"[Field]"| q85
    q85 --> q86
    q86 --> q89
    q87 --> q83
    q87 --> q88
    q88 --> q90
    q89 --> q87
    q90 -->|"tok('}')"| q91
    q91 --> q3
```

## Field

```mermaid
flowchart TD
    q4(["SN:4<br/>RuleStart"])
    q5(["SN:5<br/>RuleStop"])
    q92["SN:92<br/>Basic<br/>Terminal #1"]
    q93["SN:94<br/>Basic<br/>NonTerminal #1"]
    q94["SN:95<br/>Basic<br/>NonTerminal #1"]

    q4 --> q92
    q92 -->|"tok(ID)"| q93
    q93 -.->|"[FieldType]"| q94
    q94 --> q5
```

## FieldType

```mermaid
flowchart TD
    q6(["SN:6<br/>RuleStart"])
    q7(["SN:7<br/>RuleStop"])
    q95{"SN:95<br/>Basic<br/>Alternation #1<br/>dec=8"}
    q96["SN:96<br/>Basic<br/>NonTerminal #1"]
    q97["SN:97<br/>Basic<br/>NonTerminal #1"]
    q98["SN:98<br/>Basic<br/>NonTerminal #2"]
    q99["SN:99<br/>Basic<br/>NonTerminal #2"]
    q100["SN:100<br/>Basic<br/>NonTerminal #3"]
    q101["SN:101<br/>Basic<br/>NonTerminal #3"]
    q102["SN:102<br/>Basic<br/>NonTerminal #4"]
    q103["SN:103<br/>Basic<br/>NonTerminal #4"]
    q104["SN:104<br/>BlockEnd<br/>Alternation #1"]

    q6 --> q95
    q95 --> q96
    q95 --> q98
    q95 --> q100
    q95 --> q102
    q96 -.->|"[SimpleType]"| q97
    q97 --> q104
    q98 -.->|"[ReferenceType]"| q99
    q99 --> q104
    q100 -.->|"[ArrayType]"| q101
    q101 --> q104
    q102 -.->|"[PrimitiveType]"| q103
    q103 --> q104
    q104 --> q7
```

## ArrayType

```mermaid
flowchart TD
    q8(["SN:8<br/>RuleStart"])
    q9(["SN:9<br/>RuleStop"])
    q105["SN:105<br/>Basic<br/>Terminal #1"]
    q106["SN:107<br/>Basic<br/>Terminal #2"]
    q107["SN:109<br/>Basic<br/>NonTerminal #1"]
    q108["SN:110<br/>Basic<br/>NonTerminal #1"]

    q8 --> q105
    q105 -->|"tok('[')"| q106
    q106 -->|"tok(']')"| q107
    q107 -.->|"[FieldType]"| q108
    q108 --> q9
```

## ReferenceType

```mermaid
flowchart TD
    q10(["SN:10<br/>RuleStart"])
    q11(["SN:11<br/>RuleStop"])
    q109["SN:109<br/>Basic<br/>Terminal #1"]
    q110["SN:111<br/>Basic<br/>Terminal #2"]
    q111["SN:112<br/>Basic<br/>Terminal #2"]

    q10 --> q109
    q109 -->|"tok('*')"| q110
    q110 -->|"tok(ID)"| q111
    q111 --> q11
```

## SimpleType

```mermaid
flowchart TD
    q12(["SN:12<br/>RuleStart"])
    q13(["SN:13<br/>RuleStop"])
    q112["SN:112<br/>Basic<br/>Terminal #1"]
    q113["SN:113<br/>Basic<br/>Terminal #1"]

    q12 --> q112
    q112 -->|"tok(ID)"| q113
    q113 --> q13
```

## PrimitiveType

```mermaid
flowchart TD
    q14(["SN:14<br/>RuleStart"])
    q15(["SN:15<br/>RuleStop"])
    q114{"SN:114<br/>Basic<br/>Alternation #1<br/>dec=9"}
    q115["SN:115<br/>Basic<br/>Terminal #1"]
    q116["SN:116<br/>Basic<br/>Terminal #1"]
    q117["SN:117<br/>Basic<br/>Terminal #2"]
    q118["SN:118<br/>Basic<br/>Terminal #2"]
    q119["SN:119<br/>Basic<br/>Terminal #3"]
    q120["SN:120<br/>Basic<br/>Terminal #3"]
    q121["SN:121<br/>BlockEnd<br/>Alternation #1"]

    q14 --> q114
    q114 --> q115
    q114 --> q117
    q114 --> q119
    q115 -->|"tok('string')"| q116
    q116 --> q121
    q117 -->|"tok('bool')"| q118
    q118 --> q121
    q119 -->|"tok('composite')"| q120
    q120 --> q121
    q121 --> q15
```

## ParserRule

```mermaid
flowchart TD
    q16(["SN:16<br/>RuleStart"])
    q17(["SN:17<br/>RuleStop"])
    q122["SN:122<br/>Basic<br/>Terminal #1"]
    q123{"SN:124<br/>Basic<br/>Option #1<br/>dec=10"}
    q124["SN:125<br/>Basic<br/>Terminal #2"]
    q125["SN:127<br/>Basic<br/>Terminal #3"]
    q126["SN:128<br/>Basic<br/>Terminal #3"]
    q127["SN:128<br/>BlockEnd<br/>Option #1"]
    q128["SN:129<br/>Basic<br/>Terminal #4"]
    q129["SN:131<br/>Basic<br/>NonTerminal #1"]
    q130["SN:133<br/>Basic<br/>Terminal #5"]
    q131["SN:134<br/>Basic<br/>Terminal #5"]

    q16 --> q122
    q122 -->|"tok(ID)"| q123
    q123 --> q124
    q123 --> q127
    q124 -->|"tok('returns')"| q125
    q125 -->|"tok(ID)"| q126
    q126 --> q127
    q127 --> q128
    q128 -->|"tok(':')"| q129
    q129 -.->|"[Alternatives]"| q130
    q130 -->|"tok(';')"| q131
    q131 --> q17
```

## Token

```mermaid
flowchart TD
    q18(["SN:18<br/>RuleStart"])
    q19(["SN:19<br/>RuleStop"])
    q132{"SN:132<br/>Basic<br/>Option #1<br/>dec=11"}
    q133{"SN:133<br/>Basic<br/>Alternation #1<br/>dec=12"}
    q134["SN:134<br/>Basic<br/>Terminal #1"]
    q135["SN:135<br/>Basic<br/>Terminal #1"]
    q136["SN:136<br/>Basic<br/>Terminal #2"]
    q137["SN:137<br/>Basic<br/>Terminal #2"]
    q138["SN:138<br/>BlockEnd<br/>Alternation #1"]
    q139["SN:139<br/>BlockEnd<br/>Option #1"]
    q140["SN:140<br/>Basic<br/>Terminal #3"]
    q141["SN:142<br/>Basic<br/>Terminal #4"]
    q142["SN:144<br/>Basic<br/>Terminal #5"]
    q143["SN:146<br/>Basic<br/>Terminal #6"]
    q144["SN:148<br/>Basic<br/>Terminal #7"]
    q145["SN:149<br/>Basic<br/>Terminal #7"]

    q18 --> q132
    q132 --> q133
    q132 --> q139
    q133 --> q134
    q133 --> q136
    q134 -->|"tok('hidden')"| q135
    q135 --> q138
    q136 -->|"tok('comment')"| q137
    q137 --> q138
    q138 --> q139
    q139 --> q140
    q140 -->|"tok('token')"| q141
    q141 -->|"tok(ID)"| q142
    q142 -->|"tok(':')"| q143
    q143 -->|"tok(RegexLiteral)"| q144
    q144 -->|"tok(';')"| q145
    q145 --> q19
```

## Alternatives

```mermaid
flowchart TD
    q20(["SN:20<br/>RuleStart"])
    q21(["SN:21<br/>RuleStop"])
    q146["SN:146<br/>Basic<br/>NonTerminal #1"]
    q147{"SN:148<br/>Basic<br/>Option #1<br/>dec=13"}
    q148{"SN:149<br/>PlusBlockStart<br/>RepetitionMandatory #1<br/>dec=14"}
    q149["SN:150<br/>Basic<br/>Terminal #1"]
    q150["SN:152<br/>Basic<br/>NonTerminal #2"]
    q151["SN:153<br/>Basic<br/>NonTerminal #2"]
    q152["SN:153<br/>BlockEnd<br/>RepetitionMandatory #1"]
    q153{"SN:154<br/>PlusLoopBack<br/>RepetitionMandatory #1<br/>dec=15"}
    q154["SN:155<br/>LoopEnd<br/>RepetitionMandatory #1"]
    q155["SN:156<br/>BlockEnd<br/>Option #1"]

    q20 --> q146
    q146 -.->|"[Group]"| q147
    q147 --> q148
    q147 --> q155
    q148 --> q149
    q149 -->|"tok('|')"| q150
    q150 -.->|"[Group]"| q151
    q151 --> q152
    q152 --> q153
    q153 --> q148
    q153 --> q154
    q154 --> q155
    q155 --> q21
```

## Group

```mermaid
flowchart TD
    q22(["SN:22<br/>RuleStart"])
    q23(["SN:23<br/>RuleStop"])
    q156["SN:156<br/>Basic<br/>NonTerminal #1"]
    q157{"SN:158<br/>Basic<br/>Option #1<br/>dec=16"}
    q158{"SN:159<br/>PlusBlockStart<br/>RepetitionMandatory #1<br/>dec=17"}
    q159["SN:160<br/>Basic<br/>NonTerminal #2"]
    q160["SN:161<br/>Basic<br/>NonTerminal #2"]
    q161["SN:162<br/>BlockEnd<br/>RepetitionMandatory #1"]
    q162{"SN:163<br/>PlusLoopBack<br/>RepetitionMandatory #1<br/>dec=18"}
    q163["SN:164<br/>LoopEnd<br/>RepetitionMandatory #1"]
    q164["SN:165<br/>BlockEnd<br/>Option #1"]

    q22 --> q156
    q156 -.->|"[Element]"| q157
    q157 --> q158
    q157 --> q164
    q158 --> q159
    q159 -.->|"[Element]"| q160
    q160 --> q161
    q161 --> q162
    q162 --> q158
    q162 --> q163
    q163 --> q164
    q164 --> q23
```

## Element

```mermaid
flowchart TD
    q24(["SN:24<br/>RuleStart"])
    q25(["SN:25<br/>RuleStop"])
    q165{"SN:165<br/>Basic<br/>Alternation #1<br/>dec=19"}
    q166["SN:166<br/>Basic<br/>NonTerminal #1"]
    q167["SN:167<br/>Basic<br/>NonTerminal #1"]
    q168["SN:168<br/>Basic<br/>NonTerminal #2"]
    q169["SN:169<br/>Basic<br/>NonTerminal #2"]
    q170["SN:170<br/>Basic<br/>NonTerminal #3"]
    q171["SN:171<br/>Basic<br/>NonTerminal #3"]
    q172["SN:172<br/>Basic<br/>NonTerminal #4"]
    q173["SN:173<br/>Basic<br/>NonTerminal #4"]
    q174["SN:174<br/>Basic<br/>Terminal #1"]
    q175["SN:176<br/>Basic<br/>NonTerminal #5"]
    q176["SN:178<br/>Basic<br/>Terminal #2"]
    q177["SN:179<br/>Basic<br/>Terminal #2"]
    q178["SN:178<br/>BlockEnd<br/>Alternation #1"]
    q179{"SN:179<br/>Basic<br/>Alternation #2<br/>dec=20"}
    q180["SN:180<br/>Basic<br/>Terminal #3"]
    q181["SN:181<br/>Basic<br/>Terminal #3"]
    q182["SN:182<br/>Basic<br/>Terminal #4"]
    q183["SN:183<br/>Basic<br/>Terminal #4"]
    q184["SN:184<br/>Basic<br/>Terminal #5"]
    q185["SN:185<br/>Basic<br/>Terminal #5"]
    q186["SN:186<br/>BlockEnd<br/>Alternation #2"]

    q24 --> q165
    q165 --> q166
    q165 --> q168
    q165 --> q170
    q165 --> q172
    q165 --> q174
    q166 -.->|"[Keyword]"| q167
    q167 --> q178
    q168 -.->|"[Assignment]"| q169
    q169 --> q178
    q170 -.->|"[RuleCall]"| q171
    q171 --> q178
    q172 -.->|"[Action]"| q173
    q173 --> q178
    q174 -->|"tok('(')"| q175
    q175 -.->|"[Alternatives]"| q176
    q176 -->|"tok(')')"| q177
    q177 --> q178
    q178 --> q179
    q179 --> q180
    q179 --> q182
    q179 --> q184
    q180 -->|"tok('*')"| q181
    q181 --> q186
    q182 -->|"tok('+')"| q183
    q183 --> q186
    q184 -->|"tok('?')"| q185
    q185 --> q186
    q186 --> q25
```

## Keyword

```mermaid
flowchart TD
    q26(["SN:26<br/>RuleStart"])
    q27(["SN:27<br/>RuleStop"])
    q187["SN:187<br/>Basic<br/>Terminal #1"]
    q188["SN:188<br/>Basic<br/>Terminal #1"]

    q26 --> q187
    q187 -->|"tok(StringLiteral)"| q188
    q188 --> q27
```

## Assignment

```mermaid
flowchart TD
    q28(["SN:28<br/>RuleStart"])
    q29(["SN:29<br/>RuleStop"])
    q189["SN:189<br/>Basic<br/>Terminal #1"]
    q190{"SN:191<br/>Basic<br/>Alternation #1<br/>dec=21"}
    q191["SN:192<br/>Basic<br/>Terminal #2"]
    q192["SN:193<br/>Basic<br/>Terminal #2"]
    q193["SN:194<br/>Basic<br/>Terminal #3"]
    q194["SN:195<br/>Basic<br/>Terminal #3"]
    q195["SN:196<br/>Basic<br/>Terminal #4"]
    q196["SN:197<br/>Basic<br/>Terminal #4"]
    q197["SN:198<br/>BlockEnd<br/>Alternation #1"]
    q198["SN:199<br/>Basic<br/>NonTerminal #1"]
    q199["SN:200<br/>Basic<br/>NonTerminal #1"]

    q28 --> q189
    q189 -->|"tok(ID)"| q190
    q190 --> q191
    q190 --> q193
    q190 --> q195
    q191 -->|"tok('+=')"| q192
    q192 --> q197
    q193 -->|"tok('=')"| q194
    q194 --> q197
    q195 -->|"tok('?=')"| q196
    q196 --> q197
    q197 --> q198
    q198 -.->|"[Assignable]"| q199
    q199 --> q29
```

## Assignable

```mermaid
flowchart TD
    q30(["SN:30<br/>RuleStart"])
    q31(["SN:31<br/>RuleStop"])
    q200{"SN:200<br/>Basic<br/>Alternation #1<br/>dec=22"}
    q201["SN:201<br/>Basic<br/>NonTerminal #1"]
    q202["SN:202<br/>Basic<br/>NonTerminal #1"]
    q203["SN:203<br/>Basic<br/>NonTerminal #2"]
    q204["SN:204<br/>Basic<br/>NonTerminal #2"]
    q205["SN:205<br/>Basic<br/>NonTerminal #3"]
    q206["SN:206<br/>Basic<br/>NonTerminal #3"]
    q207["SN:207<br/>Basic<br/>Terminal #1"]
    q208["SN:209<br/>Basic<br/>NonTerminal #4"]
    q209["SN:211<br/>Basic<br/>Terminal #2"]
    q210["SN:212<br/>Basic<br/>Terminal #2"]
    q211["SN:211<br/>BlockEnd<br/>Alternation #1"]

    q30 --> q200
    q200 --> q201
    q200 --> q203
    q200 --> q205
    q200 --> q207
    q201 -.->|"[Keyword]"| q202
    q202 --> q211
    q203 -.->|"[RuleCall]"| q204
    q204 --> q211
    q205 -.->|"[CrossRef]"| q206
    q206 --> q211
    q207 -->|"tok('(')"| q208
    q208 -.->|"[AssignableAlternatives]"| q209
    q209 -->|"tok(')')"| q210
    q210 --> q211
    q211 --> q31
```

## AssignableWithoutAlts

```mermaid
flowchart TD
    q32(["SN:32<br/>RuleStart"])
    q33(["SN:33<br/>RuleStop"])
    q212{"SN:212<br/>Basic<br/>Alternation #1<br/>dec=23"}
    q213["SN:213<br/>Basic<br/>NonTerminal #1"]
    q214["SN:214<br/>Basic<br/>NonTerminal #1"]
    q215["SN:215<br/>Basic<br/>NonTerminal #2"]
    q216["SN:216<br/>Basic<br/>NonTerminal #2"]
    q217["SN:217<br/>Basic<br/>NonTerminal #3"]
    q218["SN:218<br/>Basic<br/>NonTerminal #3"]
    q219["SN:219<br/>BlockEnd<br/>Alternation #1"]

    q32 --> q212
    q212 --> q213
    q212 --> q215
    q212 --> q217
    q213 -.->|"[Keyword]"| q214
    q214 --> q219
    q215 -.->|"[RuleCall]"| q216
    q216 --> q219
    q217 -.->|"[CrossRef]"| q218
    q218 --> q219
    q219 --> q33
```

## AssignableAlternatives

```mermaid
flowchart TD
    q34(["SN:34<br/>RuleStart"])
    q35(["SN:35<br/>RuleStop"])
    q220["SN:220<br/>Basic<br/>NonTerminal #1"]
    q221{"SN:222<br/>Basic<br/>Option #1<br/>dec=24"}
    q222{"SN:223<br/>PlusBlockStart<br/>RepetitionMandatory #1<br/>dec=25"}
    q223["SN:224<br/>Basic<br/>Terminal #1"]
    q224["SN:226<br/>Basic<br/>NonTerminal #2"]
    q225["SN:227<br/>Basic<br/>NonTerminal #2"]
    q226["SN:227<br/>BlockEnd<br/>RepetitionMandatory #1"]
    q227{"SN:228<br/>PlusLoopBack<br/>RepetitionMandatory #1<br/>dec=26"}
    q228["SN:229<br/>LoopEnd<br/>RepetitionMandatory #1"]
    q229["SN:230<br/>BlockEnd<br/>Option #1"]

    q34 --> q220
    q220 -.->|"[AssignableWithoutAlts]"| q221
    q221 --> q222
    q221 --> q229
    q222 --> q223
    q223 -->|"tok('|')"| q224
    q224 -.->|"[AssignableWithoutAlts]"| q225
    q225 --> q226
    q226 --> q227
    q227 --> q222
    q227 --> q228
    q228 --> q229
    q229 --> q35
```

## CrossRef

```mermaid
flowchart TD
    q36(["SN:36<br/>RuleStart"])
    q37(["SN:37<br/>RuleStop"])
    q230["SN:230<br/>Basic<br/>Terminal #1"]
    q231["SN:232<br/>Basic<br/>Terminal #2"]
    q232{"SN:234<br/>Basic<br/>Option #1<br/>dec=27"}
    q233["SN:235<br/>Basic<br/>Terminal #3"]
    q234["SN:237<br/>Basic<br/>NonTerminal #1"]
    q235["SN:238<br/>Basic<br/>NonTerminal #1"]
    q236["SN:238<br/>BlockEnd<br/>Option #1"]
    q237["SN:239<br/>Basic<br/>Terminal #4"]
    q238["SN:240<br/>Basic<br/>Terminal #4"]

    q36 --> q230
    q230 -->|"tok('[')"| q231
    q231 -->|"tok(ID)"| q232
    q232 --> q233
    q232 --> q236
    q233 -->|"tok(':')"| q234
    q234 -.->|"[RuleCall]"| q235
    q235 --> q236
    q236 --> q237
    q237 -->|"tok(']')"| q238
    q238 --> q37
```

## RuleCall

```mermaid
flowchart TD
    q38(["SN:38<br/>RuleStart"])
    q39(["SN:39<br/>RuleStop"])
    q239["SN:239<br/>Basic<br/>Terminal #1"]
    q240["SN:240<br/>Basic<br/>Terminal #1"]

    q38 --> q239
    q239 -->|"tok(ID)"| q240
    q240 --> q39
```

## Action

```mermaid
flowchart TD
    q40(["SN:40<br/>RuleStart"])
    q41(["SN:41<br/>RuleStop"])
    q241["SN:241<br/>Basic<br/>Terminal #1"]
    q242["SN:243<br/>Basic<br/>Terminal #2"]
    q243{"SN:245<br/>Basic<br/>Option #1<br/>dec=28"}
    q244["SN:246<br/>Basic<br/>Terminal #3"]
    q245["SN:248<br/>Basic<br/>Terminal #4"]
    q246{"SN:250<br/>Basic<br/>Alternation #1<br/>dec=29"}
    q247["SN:251<br/>Basic<br/>Terminal #5"]
    q248["SN:252<br/>Basic<br/>Terminal #5"]
    q249["SN:253<br/>Basic<br/>Terminal #6"]
    q250["SN:254<br/>Basic<br/>Terminal #6"]
    q251["SN:255<br/>BlockEnd<br/>Alternation #1"]
    q252["SN:256<br/>Basic<br/>Terminal #7"]
    q253["SN:257<br/>Basic<br/>Terminal #7"]
    q254["SN:256<br/>BlockEnd<br/>Option #1"]
    q255["SN:257<br/>Basic<br/>Terminal #8"]
    q256["SN:258<br/>Basic<br/>Terminal #8"]

    q40 --> q241
    q241 -->|"tok('{')"| q242
    q242 -->|"tok(ID)"| q243
    q243 --> q244
    q243 --> q254
    q244 -->|"tok('.')"| q245
    q245 -->|"tok(ID)"| q246
    q246 --> q247
    q246 --> q249
    q247 -->|"tok('+=')"| q248
    q248 --> q251
    q249 -->|"tok('=')"| q250
    q250 --> q251
    q251 --> q252
    q252 -->|"tok('current')"| q253
    q253 --> q254
    q254 --> q255
    q255 -->|"tok('}')"| q256
    q256 --> q41
```

## CompositeRule

```mermaid
flowchart TD
    q42(["SN:42<br/>RuleStart"])
    q43(["SN:43<br/>RuleStop"])
    q257["SN:257<br/>Basic<br/>Terminal #1"]
    q258["SN:259<br/>Basic<br/>Terminal #2"]
    q259["SN:261<br/>Basic<br/>Terminal #3"]
    q260["SN:263<br/>Basic<br/>NonTerminal #1"]
    q261["SN:265<br/>Basic<br/>Terminal #4"]
    q262["SN:266<br/>Basic<br/>Terminal #4"]

    q42 --> q257
    q257 -->|"tok('composite')"| q258
    q258 -->|"tok(ID)"| q259
    q259 -->|"tok(':')"| q260
    q260 -.->|"[CompositeAlternatives]"| q261
    q261 -->|"tok(';')"| q262
    q262 --> q43
```

## CompositeAlternatives

```mermaid
flowchart TD
    q44(["SN:44<br/>RuleStart"])
    q45(["SN:45<br/>RuleStop"])
    q263["SN:263<br/>Basic<br/>NonTerminal #1"]
    q264{"SN:265<br/>Basic<br/>Option #1<br/>dec=30"}
    q265{"SN:266<br/>PlusBlockStart<br/>RepetitionMandatory #1<br/>dec=31"}
    q266["SN:267<br/>Basic<br/>Terminal #1"]
    q267["SN:269<br/>Basic<br/>NonTerminal #2"]
    q268["SN:270<br/>Basic<br/>NonTerminal #2"]
    q269["SN:270<br/>BlockEnd<br/>RepetitionMandatory #1"]
    q270{"SN:271<br/>PlusLoopBack<br/>RepetitionMandatory #1<br/>dec=32"}
    q271["SN:272<br/>LoopEnd<br/>RepetitionMandatory #1"]
    q272["SN:273<br/>BlockEnd<br/>Option #1"]

    q44 --> q263
    q263 -.->|"[CompositeGroup]"| q264
    q264 --> q265
    q264 --> q272
    q265 --> q266
    q266 -->|"tok('|')"| q267
    q267 -.->|"[CompositeGroup]"| q268
    q268 --> q269
    q269 --> q270
    q270 --> q265
    q270 --> q271
    q271 --> q272
    q272 --> q45
```

## CompositeGroup

```mermaid
flowchart TD
    q46(["SN:46<br/>RuleStart"])
    q47(["SN:47<br/>RuleStop"])
    q273["SN:273<br/>Basic<br/>NonTerminal #1"]
    q274{"SN:275<br/>Basic<br/>Option #1<br/>dec=33"}
    q275{"SN:276<br/>PlusBlockStart<br/>RepetitionMandatory #1<br/>dec=34"}
    q276["SN:277<br/>Basic<br/>NonTerminal #2"]
    q277["SN:278<br/>Basic<br/>NonTerminal #2"]
    q278["SN:279<br/>BlockEnd<br/>RepetitionMandatory #1"]
    q279{"SN:280<br/>PlusLoopBack<br/>RepetitionMandatory #1<br/>dec=35"}
    q280["SN:281<br/>LoopEnd<br/>RepetitionMandatory #1"]
    q281["SN:282<br/>BlockEnd<br/>Option #1"]

    q46 --> q273
    q273 -.->|"[CompositeElement]"| q274
    q274 --> q275
    q274 --> q281
    q275 --> q276
    q276 -.->|"[CompositeElement]"| q277
    q277 --> q278
    q278 --> q279
    q279 --> q275
    q279 --> q280
    q280 --> q281
    q281 --> q47
```

## CompositeElement

```mermaid
flowchart TD
    q48(["SN:48<br/>RuleStart"])
    q49(["SN:49<br/>RuleStop"])
    q282{"SN:282<br/>Basic<br/>Alternation #1<br/>dec=36"}
    q283["SN:283<br/>Basic<br/>NonTerminal #1"]
    q284["SN:284<br/>Basic<br/>NonTerminal #1"]
    q285["SN:285<br/>Basic<br/>NonTerminal #2"]
    q286["SN:286<br/>Basic<br/>NonTerminal #2"]
    q287["SN:287<br/>Basic<br/>Terminal #1"]
    q288["SN:289<br/>Basic<br/>NonTerminal #3"]
    q289["SN:291<br/>Basic<br/>Terminal #2"]
    q290["SN:292<br/>Basic<br/>Terminal #2"]
    q291["SN:291<br/>BlockEnd<br/>Alternation #1"]
    q292{"SN:292<br/>Basic<br/>Alternation #2<br/>dec=37"}
    q293["SN:293<br/>Basic<br/>Terminal #3"]
    q294["SN:294<br/>Basic<br/>Terminal #3"]
    q295["SN:295<br/>Basic<br/>Terminal #4"]
    q296["SN:296<br/>Basic<br/>Terminal #4"]
    q297["SN:297<br/>Basic<br/>Terminal #5"]
    q298["SN:298<br/>Basic<br/>Terminal #5"]
    q299["SN:299<br/>BlockEnd<br/>Alternation #2"]

    q48 --> q282
    q282 --> q283
    q282 --> q285
    q282 --> q287
    q283 -.->|"[Keyword]"| q284
    q284 --> q291
    q285 -.->|"[RuleCall]"| q286
    q286 --> q291
    q287 -->|"tok('(')"| q288
    q288 -.->|"[CompositeAlternatives]"| q289
    q289 -->|"tok(')')"| q290
    q290 --> q291
    q291 --> q292
    q292 --> q293
    q292 --> q295
    q292 --> q297
    q293 -->|"tok('*')"| q294
    q294 --> q299
    q295 -->|"tok('+')"| q296
    q296 --> q299
    q297 -->|"tok('?')"| q298
    q298 --> q299
    q299 --> q49
```

