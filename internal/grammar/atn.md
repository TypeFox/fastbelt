# Runtime ATN for grammar

## Grammar

```mermaid
flowchart TD
    q0(["SN:0<br/>RuleStart"])
    q1(["SN:1<br/>RuleStop"])
    q50["SN:50<br/>Basic<br/>Terminal #1"]
    q51["SN:52<br/>Basic<br/>Terminal #2"]
    q52["SN:54<br/>Basic<br/>Terminal #3"]
    q53{"SN:56<br/>Basic<br/>Alternation #1<br/>dec=0"}
    q54["SN:57<br/>Basic<br/>NonTerminal #1"]
    q55["SN:58<br/>Basic<br/>NonTerminal #1"]
    q56["SN:59<br/>Basic<br/>NonTerminal #2"]
    q57["SN:60<br/>Basic<br/>NonTerminal #2"]
    q58["SN:61<br/>Basic<br/>NonTerminal #3"]
    q59["SN:62<br/>Basic<br/>NonTerminal #3"]
    q60["SN:63<br/>Basic<br/>NonTerminal #4"]
    q61["SN:64<br/>Basic<br/>NonTerminal #4"]
    q62["SN:65<br/>BlockEnd<br/>Alternation #1"]

    q0 --> q50
    q50 -->|"tok('grammar')"| q51
    q51 -->|"tok(ID)"| q52
    q52 -->|"tok(';')"| q53
    q53 --> q54
    q53 --> q56
    q53 --> q58
    q53 --> q60
    q54 -.->|"[ParserRule]"| q55
    q55 --> q62
    q56 -.->|"[Token]"| q57
    q57 --> q62
    q58 -.->|"[Interface]"| q59
    q59 --> q62
    q60 -.->|"[CompositeRule]"| q61
    q61 --> q62
    q62 --> q1
```

## Interface

```mermaid
flowchart TD
    q2(["SN:2<br/>RuleStart"])
    q3(["SN:3<br/>RuleStop"])
    q63["SN:63<br/>Basic<br/>Terminal #1"]
    q64["SN:65<br/>Basic<br/>Terminal #2"]
    q65{"SN:67<br/>Basic<br/>Option #1<br/>dec=1"}
    q66["SN:68<br/>Basic<br/>Terminal #3"]
    q67["SN:70<br/>Basic<br/>Terminal #4"]
    q68{"SN:72<br/>StarBlockStart<br/>Repetition #1<br/>dec=2"}
    q69["SN:73<br/>Basic<br/>Terminal #5"]
    q70["SN:75<br/>Basic<br/>Terminal #6"]
    q71["SN:76<br/>Basic<br/>Terminal #6"]
    q72["SN:76<br/>BlockEnd<br/>Repetition #1"]
    q73{"SN:77<br/>StarLoopEntry<br/>Repetition #1<br/>dec=3"}
    q74["SN:78<br/>LoopEnd<br/>Repetition #1"]
    q75["SN:79<br/>StarLoopBack<br/>Repetition #1"]
    q76["SN:78<br/>BlockEnd<br/>Option #1"]
    q77["SN:79<br/>Basic<br/>Terminal #7"]
    q78{"SN:81<br/>StarBlockStart<br/>Repetition #2<br/>dec=4"}
    q79["SN:82<br/>Basic<br/>NonTerminal #1"]
    q80["SN:83<br/>Basic<br/>NonTerminal #1"]
    q81["SN:84<br/>BlockEnd<br/>Repetition #2"]
    q82{"SN:85<br/>StarLoopEntry<br/>Repetition #2<br/>dec=5"}
    q83["SN:86<br/>LoopEnd<br/>Repetition #2"]
    q84["SN:87<br/>StarLoopBack<br/>Repetition #2"]
    q85["SN:88<br/>Basic<br/>Terminal #8"]
    q86["SN:89<br/>Basic<br/>Terminal #8"]

    q2 --> q63
    q63 -->|"tok('interface')"| q64
    q64 -->|"tok(ID)"| q65
    q65 --> q66
    q65 --> q76
    q66 -->|"tok('extends')"| q67
    q67 -->|"tok(ID)"| q73
    q68 --> q69
    q69 -->|"tok(',')"| q70
    q70 -->|"tok(ID)"| q71
    q71 --> q72
    q72 --> q75
    q73 --> q68
    q73 --> q74
    q74 --> q76
    q75 --> q73
    q76 --> q77
    q77 -->|"tok('{')"| q82
    q78 --> q79
    q79 -.->|"[Field]"| q80
    q80 --> q81
    q81 --> q84
    q82 --> q78
    q82 --> q83
    q83 --> q85
    q84 --> q82
    q85 -->|"tok('}')"| q86
    q86 --> q3
```

## Field

```mermaid
flowchart TD
    q4(["SN:4<br/>RuleStart"])
    q5(["SN:5<br/>RuleStop"])
    q87["SN:87<br/>Basic<br/>Terminal #1"]
    q88["SN:89<br/>Basic<br/>NonTerminal #1"]
    q89["SN:90<br/>Basic<br/>NonTerminal #1"]

    q4 --> q87
    q87 -->|"tok(ID)"| q88
    q88 -.->|"[FieldType]"| q89
    q89 --> q5
```

## FieldType

```mermaid
flowchart TD
    q6(["SN:6<br/>RuleStart"])
    q7(["SN:7<br/>RuleStop"])
    q90{"SN:90<br/>Basic<br/>Alternation #1<br/>dec=6"}
    q91["SN:91<br/>Basic<br/>NonTerminal #1"]
    q92["SN:92<br/>Basic<br/>NonTerminal #1"]
    q93["SN:93<br/>Basic<br/>NonTerminal #2"]
    q94["SN:94<br/>Basic<br/>NonTerminal #2"]
    q95["SN:95<br/>Basic<br/>NonTerminal #3"]
    q96["SN:96<br/>Basic<br/>NonTerminal #3"]
    q97["SN:97<br/>Basic<br/>NonTerminal #4"]
    q98["SN:98<br/>Basic<br/>NonTerminal #4"]
    q99["SN:99<br/>BlockEnd<br/>Alternation #1"]

    q6 --> q90
    q90 --> q91
    q90 --> q93
    q90 --> q95
    q90 --> q97
    q91 -.->|"[SimpleType]"| q92
    q92 --> q99
    q93 -.->|"[ReferenceType]"| q94
    q94 --> q99
    q95 -.->|"[ArrayType]"| q96
    q96 --> q99
    q97 -.->|"[PrimitiveType]"| q98
    q98 --> q99
    q99 --> q7
```

## ArrayType

```mermaid
flowchart TD
    q8(["SN:8<br/>RuleStart"])
    q9(["SN:9<br/>RuleStop"])
    q100["SN:100<br/>Basic<br/>Terminal #1"]
    q101["SN:102<br/>Basic<br/>Terminal #2"]
    q102["SN:104<br/>Basic<br/>NonTerminal #1"]
    q103["SN:105<br/>Basic<br/>NonTerminal #1"]

    q8 --> q100
    q100 -->|"tok('[')"| q101
    q101 -->|"tok(']')"| q102
    q102 -.->|"[FieldType]"| q103
    q103 --> q9
```

## ReferenceType

```mermaid
flowchart TD
    q10(["SN:10<br/>RuleStart"])
    q11(["SN:11<br/>RuleStop"])
    q104["SN:104<br/>Basic<br/>Terminal #1"]
    q105["SN:106<br/>Basic<br/>Terminal #2"]
    q106["SN:107<br/>Basic<br/>Terminal #2"]

    q10 --> q104
    q104 -->|"tok('*')"| q105
    q105 -->|"tok(ID)"| q106
    q106 --> q11
```

## SimpleType

```mermaid
flowchart TD
    q12(["SN:12<br/>RuleStart"])
    q13(["SN:13<br/>RuleStop"])
    q107["SN:107<br/>Basic<br/>Terminal #1"]
    q108["SN:108<br/>Basic<br/>Terminal #1"]

    q12 --> q107
    q107 -->|"tok(ID)"| q108
    q108 --> q13
```

## PrimitiveType

```mermaid
flowchart TD
    q14(["SN:14<br/>RuleStart"])
    q15(["SN:15<br/>RuleStop"])
    q109{"SN:109<br/>Basic<br/>Alternation #1<br/>dec=7"}
    q110["SN:110<br/>Basic<br/>Terminal #1"]
    q111["SN:111<br/>Basic<br/>Terminal #1"]
    q112["SN:112<br/>Basic<br/>Terminal #2"]
    q113["SN:113<br/>Basic<br/>Terminal #2"]
    q114["SN:114<br/>Basic<br/>Terminal #3"]
    q115["SN:115<br/>Basic<br/>Terminal #3"]
    q116["SN:116<br/>BlockEnd<br/>Alternation #1"]

    q14 --> q109
    q109 --> q110
    q109 --> q112
    q109 --> q114
    q110 -->|"tok('string')"| q111
    q111 --> q116
    q112 -->|"tok('bool')"| q113
    q113 --> q116
    q114 -->|"tok('composite')"| q115
    q115 --> q116
    q116 --> q15
```

## ParserRule

```mermaid
flowchart TD
    q16(["SN:16<br/>RuleStart"])
    q17(["SN:17<br/>RuleStop"])
    q117["SN:117<br/>Basic<br/>Terminal #1"]
    q118{"SN:119<br/>Basic<br/>Option #1<br/>dec=8"}
    q119["SN:120<br/>Basic<br/>Terminal #2"]
    q120["SN:122<br/>Basic<br/>Terminal #3"]
    q121["SN:123<br/>Basic<br/>Terminal #3"]
    q122["SN:123<br/>BlockEnd<br/>Option #1"]
    q123["SN:124<br/>Basic<br/>Terminal #4"]
    q124["SN:126<br/>Basic<br/>NonTerminal #1"]
    q125["SN:128<br/>Basic<br/>Terminal #5"]
    q126["SN:129<br/>Basic<br/>Terminal #5"]

    q16 --> q117
    q117 -->|"tok(ID)"| q118
    q118 --> q119
    q118 --> q122
    q119 -->|"tok('returns')"| q120
    q120 -->|"tok(ID)"| q121
    q121 --> q122
    q122 --> q123
    q123 -->|"tok(':')"| q124
    q124 -.->|"[Alternatives]"| q125
    q125 -->|"tok(';')"| q126
    q126 --> q17
```

## Token

```mermaid
flowchart TD
    q18(["SN:18<br/>RuleStart"])
    q19(["SN:19<br/>RuleStop"])
    q127{"SN:127<br/>Basic<br/>Alternation #1<br/>dec=9"}
    q128["SN:128<br/>Basic<br/>Terminal #1"]
    q129["SN:129<br/>Basic<br/>Terminal #1"]
    q130["SN:130<br/>Basic<br/>Terminal #2"]
    q131["SN:131<br/>Basic<br/>Terminal #2"]
    q132["SN:132<br/>BlockEnd<br/>Alternation #1"]
    q133["SN:133<br/>Basic<br/>Terminal #3"]
    q134["SN:135<br/>Basic<br/>Terminal #4"]
    q135["SN:137<br/>Basic<br/>Terminal #5"]
    q136["SN:139<br/>Basic<br/>Terminal #6"]
    q137["SN:141<br/>Basic<br/>Terminal #7"]
    q138["SN:142<br/>Basic<br/>Terminal #7"]

    q18 --> q127
    q127 --> q128
    q127 --> q130
    q128 -->|"tok('hidden')"| q129
    q129 --> q132
    q130 -->|"tok('comment')"| q131
    q131 --> q132
    q132 --> q133
    q133 -->|"tok('token')"| q134
    q134 -->|"tok(ID)"| q135
    q135 -->|"tok(':')"| q136
    q136 -->|"tok(RegexLiteral)"| q137
    q137 -->|"tok(';')"| q138
    q138 --> q19
```

## Alternatives

```mermaid
flowchart TD
    q20(["SN:20<br/>RuleStart"])
    q21(["SN:21<br/>RuleStop"])
    q139["SN:139<br/>Basic<br/>NonTerminal #1"]
    q140{"SN:141<br/>Basic<br/>Option #1<br/>dec=10"}
    q141{"SN:142<br/>PlusBlockStart<br/>RepetitionMandatory #1<br/>dec=11"}
    q142["SN:143<br/>Basic<br/>Terminal #1"]
    q143["SN:145<br/>Basic<br/>NonTerminal #2"]
    q144["SN:146<br/>Basic<br/>NonTerminal #2"]
    q145["SN:146<br/>BlockEnd<br/>RepetitionMandatory #1"]
    q146{"SN:147<br/>PlusLoopBack<br/>RepetitionMandatory #1<br/>dec=12"}
    q147["SN:148<br/>LoopEnd<br/>RepetitionMandatory #1"]
    q148["SN:149<br/>BlockEnd<br/>Option #1"]

    q20 --> q139
    q139 -.->|"[Group]"| q140
    q140 --> q141
    q140 --> q148
    q141 --> q142
    q142 -->|"tok('|')"| q143
    q143 -.->|"[Group]"| q144
    q144 --> q145
    q145 --> q146
    q146 --> q141
    q146 --> q147
    q147 --> q148
    q148 --> q21
```

## Group

```mermaid
flowchart TD
    q22(["SN:22<br/>RuleStart"])
    q23(["SN:23<br/>RuleStop"])
    q149["SN:149<br/>Basic<br/>NonTerminal #1"]
    q150{"SN:151<br/>Basic<br/>Option #1<br/>dec=13"}
    q151{"SN:152<br/>PlusBlockStart<br/>RepetitionMandatory #1<br/>dec=14"}
    q152["SN:153<br/>Basic<br/>NonTerminal #2"]
    q153["SN:154<br/>Basic<br/>NonTerminal #2"]
    q154["SN:155<br/>BlockEnd<br/>RepetitionMandatory #1"]
    q155{"SN:156<br/>PlusLoopBack<br/>RepetitionMandatory #1<br/>dec=15"}
    q156["SN:157<br/>LoopEnd<br/>RepetitionMandatory #1"]
    q157["SN:158<br/>BlockEnd<br/>Option #1"]

    q22 --> q149
    q149 -.->|"[Element]"| q150
    q150 --> q151
    q150 --> q157
    q151 --> q152
    q152 -.->|"[Element]"| q153
    q153 --> q154
    q154 --> q155
    q155 --> q151
    q155 --> q156
    q156 --> q157
    q157 --> q23
```

## Element

```mermaid
flowchart TD
    q24(["SN:24<br/>RuleStart"])
    q25(["SN:25<br/>RuleStop"])
    q158{"SN:158<br/>Basic<br/>Alternation #1<br/>dec=16"}
    q159["SN:159<br/>Basic<br/>NonTerminal #1"]
    q160["SN:160<br/>Basic<br/>NonTerminal #1"]
    q161["SN:161<br/>Basic<br/>NonTerminal #2"]
    q162["SN:162<br/>Basic<br/>NonTerminal #2"]
    q163["SN:163<br/>Basic<br/>NonTerminal #3"]
    q164["SN:164<br/>Basic<br/>NonTerminal #3"]
    q165["SN:165<br/>Basic<br/>NonTerminal #4"]
    q166["SN:166<br/>Basic<br/>NonTerminal #4"]
    q167["SN:167<br/>Basic<br/>Terminal #1"]
    q168["SN:169<br/>Basic<br/>NonTerminal #5"]
    q169["SN:171<br/>Basic<br/>Terminal #2"]
    q170["SN:172<br/>Basic<br/>Terminal #2"]
    q171["SN:171<br/>BlockEnd<br/>Alternation #1"]
    q172{"SN:172<br/>Basic<br/>Alternation #2<br/>dec=17"}
    q173["SN:173<br/>Basic<br/>Terminal #3"]
    q174["SN:174<br/>Basic<br/>Terminal #3"]
    q175["SN:175<br/>Basic<br/>Terminal #4"]
    q176["SN:176<br/>Basic<br/>Terminal #4"]
    q177["SN:177<br/>Basic<br/>Terminal #5"]
    q178["SN:178<br/>Basic<br/>Terminal #5"]
    q179["SN:179<br/>BlockEnd<br/>Alternation #2"]

    q24 --> q158
    q158 --> q159
    q158 --> q161
    q158 --> q163
    q158 --> q165
    q158 --> q167
    q159 -.->|"[Keyword]"| q160
    q160 --> q171
    q161 -.->|"[Assignment]"| q162
    q162 --> q171
    q163 -.->|"[RuleCall]"| q164
    q164 --> q171
    q165 -.->|"[Action]"| q166
    q166 --> q171
    q167 -->|"tok('(')"| q168
    q168 -.->|"[Alternatives]"| q169
    q169 -->|"tok(')')"| q170
    q170 --> q171
    q171 --> q172
    q172 --> q173
    q172 --> q175
    q172 --> q177
    q173 -->|"tok('*')"| q174
    q174 --> q179
    q175 -->|"tok('+')"| q176
    q176 --> q179
    q177 -->|"tok('?')"| q178
    q178 --> q179
    q179 --> q25
```

## Keyword

```mermaid
flowchart TD
    q26(["SN:26<br/>RuleStart"])
    q27(["SN:27<br/>RuleStop"])
    q180["SN:180<br/>Basic<br/>Terminal #1"]
    q181["SN:181<br/>Basic<br/>Terminal #1"]

    q26 --> q180
    q180 -->|"tok(StringLiteral)"| q181
    q181 --> q27
```

## Assignment

```mermaid
flowchart TD
    q28(["SN:28<br/>RuleStart"])
    q29(["SN:29<br/>RuleStop"])
    q182["SN:182<br/>Basic<br/>Terminal #1"]
    q183{"SN:184<br/>Basic<br/>Alternation #1<br/>dec=18"}
    q184["SN:185<br/>Basic<br/>Terminal #2"]
    q185["SN:186<br/>Basic<br/>Terminal #2"]
    q186["SN:187<br/>Basic<br/>Terminal #3"]
    q187["SN:188<br/>Basic<br/>Terminal #3"]
    q188["SN:189<br/>Basic<br/>Terminal #4"]
    q189["SN:190<br/>Basic<br/>Terminal #4"]
    q190["SN:191<br/>BlockEnd<br/>Alternation #1"]
    q191["SN:192<br/>Basic<br/>NonTerminal #1"]
    q192["SN:193<br/>Basic<br/>NonTerminal #1"]

    q28 --> q182
    q182 -->|"tok(ID)"| q183
    q183 --> q184
    q183 --> q186
    q183 --> q188
    q184 -->|"tok('+=')"| q185
    q185 --> q190
    q186 -->|"tok('=')"| q187
    q187 --> q190
    q188 -->|"tok('?=')"| q189
    q189 --> q190
    q190 --> q191
    q191 -.->|"[Assignable]"| q192
    q192 --> q29
```

## Assignable

```mermaid
flowchart TD
    q30(["SN:30<br/>RuleStart"])
    q31(["SN:31<br/>RuleStop"])
    q193{"SN:193<br/>Basic<br/>Alternation #1<br/>dec=19"}
    q194["SN:194<br/>Basic<br/>NonTerminal #1"]
    q195["SN:195<br/>Basic<br/>NonTerminal #1"]
    q196["SN:196<br/>Basic<br/>NonTerminal #2"]
    q197["SN:197<br/>Basic<br/>NonTerminal #2"]
    q198["SN:198<br/>Basic<br/>NonTerminal #3"]
    q199["SN:199<br/>Basic<br/>NonTerminal #3"]
    q200["SN:200<br/>Basic<br/>Terminal #1"]
    q201["SN:202<br/>Basic<br/>NonTerminal #4"]
    q202["SN:204<br/>Basic<br/>Terminal #2"]
    q203["SN:205<br/>Basic<br/>Terminal #2"]
    q204["SN:204<br/>BlockEnd<br/>Alternation #1"]

    q30 --> q193
    q193 --> q194
    q193 --> q196
    q193 --> q198
    q193 --> q200
    q194 -.->|"[Keyword]"| q195
    q195 --> q204
    q196 -.->|"[RuleCall]"| q197
    q197 --> q204
    q198 -.->|"[CrossRef]"| q199
    q199 --> q204
    q200 -->|"tok('(')"| q201
    q201 -.->|"[AssignableAlternatives]"| q202
    q202 -->|"tok(')')"| q203
    q203 --> q204
    q204 --> q31
```

## AssignableWithoutAlts

```mermaid
flowchart TD
    q32(["SN:32<br/>RuleStart"])
    q33(["SN:33<br/>RuleStop"])
    q205{"SN:205<br/>Basic<br/>Alternation #1<br/>dec=20"}
    q206["SN:206<br/>Basic<br/>NonTerminal #1"]
    q207["SN:207<br/>Basic<br/>NonTerminal #1"]
    q208["SN:208<br/>Basic<br/>NonTerminal #2"]
    q209["SN:209<br/>Basic<br/>NonTerminal #2"]
    q210["SN:210<br/>Basic<br/>NonTerminal #3"]
    q211["SN:211<br/>Basic<br/>NonTerminal #3"]
    q212["SN:212<br/>BlockEnd<br/>Alternation #1"]

    q32 --> q205
    q205 --> q206
    q205 --> q208
    q205 --> q210
    q206 -.->|"[Keyword]"| q207
    q207 --> q212
    q208 -.->|"[RuleCall]"| q209
    q209 --> q212
    q210 -.->|"[CrossRef]"| q211
    q211 --> q212
    q212 --> q33
```

## AssignableAlternatives

```mermaid
flowchart TD
    q34(["SN:34<br/>RuleStart"])
    q35(["SN:35<br/>RuleStop"])
    q213["SN:213<br/>Basic<br/>NonTerminal #1"]
    q214{"SN:215<br/>Basic<br/>Option #1<br/>dec=21"}
    q215{"SN:216<br/>PlusBlockStart<br/>RepetitionMandatory #1<br/>dec=22"}
    q216["SN:217<br/>Basic<br/>Terminal #1"]
    q217["SN:219<br/>Basic<br/>NonTerminal #2"]
    q218["SN:220<br/>Basic<br/>NonTerminal #2"]
    q219["SN:220<br/>BlockEnd<br/>RepetitionMandatory #1"]
    q220{"SN:221<br/>PlusLoopBack<br/>RepetitionMandatory #1<br/>dec=23"}
    q221["SN:222<br/>LoopEnd<br/>RepetitionMandatory #1"]
    q222["SN:223<br/>BlockEnd<br/>Option #1"]

    q34 --> q213
    q213 -.->|"[AssignableWithoutAlts]"| q214
    q214 --> q215
    q214 --> q222
    q215 --> q216
    q216 -->|"tok('|')"| q217
    q217 -.->|"[AssignableWithoutAlts]"| q218
    q218 --> q219
    q219 --> q220
    q220 --> q215
    q220 --> q221
    q221 --> q222
    q222 --> q35
```

## CrossRef

```mermaid
flowchart TD
    q36(["SN:36<br/>RuleStart"])
    q37(["SN:37<br/>RuleStop"])
    q223["SN:223<br/>Basic<br/>Terminal #1"]
    q224["SN:225<br/>Basic<br/>Terminal #2"]
    q225{"SN:227<br/>Basic<br/>Option #1<br/>dec=24"}
    q226["SN:228<br/>Basic<br/>Terminal #3"]
    q227["SN:230<br/>Basic<br/>NonTerminal #1"]
    q228["SN:231<br/>Basic<br/>NonTerminal #1"]
    q229["SN:231<br/>BlockEnd<br/>Option #1"]
    q230["SN:232<br/>Basic<br/>Terminal #4"]
    q231["SN:233<br/>Basic<br/>Terminal #4"]

    q36 --> q223
    q223 -->|"tok('[')"| q224
    q224 -->|"tok(ID)"| q225
    q225 --> q226
    q225 --> q229
    q226 -->|"tok(':')"| q227
    q227 -.->|"[RuleCall]"| q228
    q228 --> q229
    q229 --> q230
    q230 -->|"tok(']')"| q231
    q231 --> q37
```

## RuleCall

```mermaid
flowchart TD
    q38(["SN:38<br/>RuleStart"])
    q39(["SN:39<br/>RuleStop"])
    q232["SN:232<br/>Basic<br/>Terminal #1"]
    q233["SN:233<br/>Basic<br/>Terminal #1"]

    q38 --> q232
    q232 -->|"tok(ID)"| q233
    q233 --> q39
```

## Action

```mermaid
flowchart TD
    q40(["SN:40<br/>RuleStart"])
    q41(["SN:41<br/>RuleStop"])
    q234["SN:234<br/>Basic<br/>Terminal #1"]
    q235["SN:236<br/>Basic<br/>Terminal #2"]
    q236{"SN:238<br/>Basic<br/>Option #1<br/>dec=25"}
    q237["SN:239<br/>Basic<br/>Terminal #3"]
    q238["SN:241<br/>Basic<br/>Terminal #4"]
    q239{"SN:243<br/>Basic<br/>Alternation #1<br/>dec=26"}
    q240["SN:244<br/>Basic<br/>Terminal #5"]
    q241["SN:245<br/>Basic<br/>Terminal #5"]
    q242["SN:246<br/>Basic<br/>Terminal #6"]
    q243["SN:247<br/>Basic<br/>Terminal #6"]
    q244["SN:248<br/>BlockEnd<br/>Alternation #1"]
    q245["SN:249<br/>Basic<br/>Terminal #7"]
    q246["SN:250<br/>Basic<br/>Terminal #7"]
    q247["SN:249<br/>BlockEnd<br/>Option #1"]
    q248["SN:250<br/>Basic<br/>Terminal #8"]
    q249["SN:251<br/>Basic<br/>Terminal #8"]

    q40 --> q234
    q234 -->|"tok('{')"| q235
    q235 -->|"tok(ID)"| q236
    q236 --> q237
    q236 --> q247
    q237 -->|"tok('.')"| q238
    q238 -->|"tok(ID)"| q239
    q239 --> q240
    q239 --> q242
    q240 -->|"tok('+=')"| q241
    q241 --> q244
    q242 -->|"tok('=')"| q243
    q243 --> q244
    q244 --> q245
    q245 -->|"tok('current')"| q246
    q246 --> q247
    q247 --> q248
    q248 -->|"tok('}')"| q249
    q249 --> q41
```

## CompositeRule

```mermaid
flowchart TD
    q42(["SN:42<br/>RuleStart"])
    q43(["SN:43<br/>RuleStop"])
    q250["SN:250<br/>Basic<br/>Terminal #1"]
    q251["SN:252<br/>Basic<br/>Terminal #2"]
    q252["SN:254<br/>Basic<br/>Terminal #3"]
    q253["SN:256<br/>Basic<br/>NonTerminal #1"]
    q254["SN:258<br/>Basic<br/>Terminal #4"]
    q255["SN:259<br/>Basic<br/>Terminal #4"]

    q42 --> q250
    q250 -->|"tok('composite')"| q251
    q251 -->|"tok(ID)"| q252
    q252 -->|"tok(':')"| q253
    q253 -.->|"[CompositeAlternatives]"| q254
    q254 -->|"tok(';')"| q255
    q255 --> q43
```

## CompositeAlternatives

```mermaid
flowchart TD
    q44(["SN:44<br/>RuleStart"])
    q45(["SN:45<br/>RuleStop"])
    q256["SN:256<br/>Basic<br/>NonTerminal #1"]
    q257{"SN:258<br/>Basic<br/>Option #1<br/>dec=27"}
    q258{"SN:259<br/>PlusBlockStart<br/>RepetitionMandatory #1<br/>dec=28"}
    q259["SN:260<br/>Basic<br/>Terminal #1"]
    q260["SN:262<br/>Basic<br/>NonTerminal #2"]
    q261["SN:263<br/>Basic<br/>NonTerminal #2"]
    q262["SN:263<br/>BlockEnd<br/>RepetitionMandatory #1"]
    q263{"SN:264<br/>PlusLoopBack<br/>RepetitionMandatory #1<br/>dec=29"}
    q264["SN:265<br/>LoopEnd<br/>RepetitionMandatory #1"]
    q265["SN:266<br/>BlockEnd<br/>Option #1"]

    q44 --> q256
    q256 -.->|"[CompositeGroup]"| q257
    q257 --> q258
    q257 --> q265
    q258 --> q259
    q259 -->|"tok('|')"| q260
    q260 -.->|"[CompositeGroup]"| q261
    q261 --> q262
    q262 --> q263
    q263 --> q258
    q263 --> q264
    q264 --> q265
    q265 --> q45
```

## CompositeGroup

```mermaid
flowchart TD
    q46(["SN:46<br/>RuleStart"])
    q47(["SN:47<br/>RuleStop"])
    q266["SN:266<br/>Basic<br/>NonTerminal #1"]
    q267{"SN:268<br/>Basic<br/>Option #1<br/>dec=30"}
    q268{"SN:269<br/>PlusBlockStart<br/>RepetitionMandatory #1<br/>dec=31"}
    q269["SN:270<br/>Basic<br/>NonTerminal #2"]
    q270["SN:271<br/>Basic<br/>NonTerminal #2"]
    q271["SN:272<br/>BlockEnd<br/>RepetitionMandatory #1"]
    q272{"SN:273<br/>PlusLoopBack<br/>RepetitionMandatory #1<br/>dec=32"}
    q273["SN:274<br/>LoopEnd<br/>RepetitionMandatory #1"]
    q274["SN:275<br/>BlockEnd<br/>Option #1"]

    q46 --> q266
    q266 -.->|"[CompositeElement]"| q267
    q267 --> q268
    q267 --> q274
    q268 --> q269
    q269 -.->|"[CompositeElement]"| q270
    q270 --> q271
    q271 --> q272
    q272 --> q268
    q272 --> q273
    q273 --> q274
    q274 --> q47
```

## CompositeElement

```mermaid
flowchart TD
    q48(["SN:48<br/>RuleStart"])
    q49(["SN:49<br/>RuleStop"])
    q275{"SN:275<br/>Basic<br/>Alternation #1<br/>dec=33"}
    q276["SN:276<br/>Basic<br/>NonTerminal #1"]
    q277["SN:277<br/>Basic<br/>NonTerminal #1"]
    q278["SN:278<br/>Basic<br/>NonTerminal #2"]
    q279["SN:279<br/>Basic<br/>NonTerminal #2"]
    q280["SN:280<br/>Basic<br/>Terminal #1"]
    q281["SN:282<br/>Basic<br/>NonTerminal #3"]
    q282["SN:284<br/>Basic<br/>Terminal #2"]
    q283["SN:285<br/>Basic<br/>Terminal #2"]
    q284["SN:284<br/>BlockEnd<br/>Alternation #1"]
    q285{"SN:285<br/>Basic<br/>Alternation #2<br/>dec=34"}
    q286["SN:286<br/>Basic<br/>Terminal #3"]
    q287["SN:287<br/>Basic<br/>Terminal #3"]
    q288["SN:288<br/>Basic<br/>Terminal #4"]
    q289["SN:289<br/>Basic<br/>Terminal #4"]
    q290["SN:290<br/>Basic<br/>Terminal #5"]
    q291["SN:291<br/>Basic<br/>Terminal #5"]
    q292["SN:292<br/>BlockEnd<br/>Alternation #2"]

    q48 --> q275
    q275 --> q276
    q275 --> q278
    q275 --> q280
    q276 -.->|"[Keyword]"| q277
    q277 --> q284
    q278 -.->|"[RuleCall]"| q279
    q279 --> q284
    q280 -->|"tok('(')"| q281
    q281 -.->|"[CompositeAlternatives]"| q282
    q282 -->|"tok(')')"| q283
    q283 --> q284
    q284 --> q285
    q285 --> q286
    q285 --> q288
    q285 --> q290
    q286 -->|"tok('*')"| q287
    q287 --> q292
    q288 -->|"tok('+')"| q289
    q289 --> q292
    q290 -->|"tok('?')"| q291
    q291 --> q292
    q292 --> q49
```

