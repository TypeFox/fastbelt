# Runtime ATN for grammar

## Grammar

```mermaid
flowchart TD
    q0(["Grammar__Start (0)<br/>RuleStart"])
    q1(["Grammar__Stop (1)<br/>RuleStop"])
    q56["Grammar_grammar (56)<br/>Basic<br/>"]
    q57["Grammar_Name_ID (57)<br/>Basic<br/>"]
    q58["Grammar_Semicolon (58)<br/>Basic<br/>"]
    q59["Grammar__Basic_0 (59)<br/>Basic<br/>"]
    q60{"Grammar__Basic_1 (60)<br/>Basic<br/><br/>dec=0"}
    q61["Grammar__Basic_2 (61)<br/>Basic<br/>"]
    q62["Grammar__Basic_3 (62)<br/>Basic<br/>"]
    q63["Grammar__Basic_4 (63)<br/>Basic<br/>"]
    q64["Grammar__Basic_5 (64)<br/>Basic<br/>"]
    q65["Grammar__Basic_6 (65)<br/>Basic<br/>"]
    q66["Grammar__Basic_7 (66)<br/>Basic<br/>"]
    q67["Grammar__Basic_8 (67)<br/>Basic<br/>"]
    q68["Grammar__Basic_9 (68)<br/>Basic<br/>"]
    q69["Grammar__Basic_10 (69)<br/>Basic<br/>"]
    q70["Grammar__Basic_11 (70)<br/>Basic<br/>"]
    q71["Grammar__Basic_12 (71)<br/>Basic<br/>"]
    q72["Grammar__Basic_13 (72)<br/>Basic<br/>"]
    q73{"Grammar__Basic_14 (73)<br/>Basic<br/><br/>dec=1"}
    q74["Grammar__BlockEnd (74)<br/>BlockEnd<br/>"]
    q75{"Grammar__LoopEntry (75)<br/>LoopEntry<br/><br/>dec=2"}
    q76["Grammar__LoopEnd (76)<br/>LoopEnd<br/>"]
    q77["Grammar__LoopBack (77)<br/>LoopBack<br/>"]

    q0 --> q56
    q56 -->|"tok(&quot;grammar&quot;)"| q57
    q57 -->|"tok(ID)"| q60
    q58 -->|"tok(&quot;;&quot;)"| q59
    q59 --> q75
    q60 --> q58
    q60 --> q59
    q61 -.->|"[ParserRule]"| q62
    q62 --> q74
    q63 -.->|"[Token]"| q64
    q64 --> q74
    q65 -.->|"[TokenGroup]"| q66
    q66 --> q74
    q67 -.->|"[TokenMode]"| q68
    q68 --> q74
    q69 -.->|"[Interface]"| q70
    q70 --> q74
    q71 -.->|"[CompositeRule]"| q72
    q72 --> q74
    q73 --> q61
    q73 --> q63
    q73 --> q65
    q73 --> q67
    q73 --> q69
    q73 --> q71
    q74 --> q77
    q75 --> q73
    q75 --> q76
    q76 --> q1
    q77 --> q75
```

## Interface

```mermaid
flowchart TD
    q2(["Interface__Start (2)<br/>RuleStart"])
    q3(["Interface__Stop (3)<br/>RuleStop"])
    q78["Interface_interface (78)<br/>Basic<br/>"]
    q79["Interface_Name_ID (79)<br/>Basic<br/>"]
    q80["Interface_extends (80)<br/>Basic<br/>"]
    q81["Interface_Extends_ID_0 (81)<br/>Basic<br/>"]
    q82["Interface_Comma (82)<br/>Basic<br/>"]
    q83["Interface_Extends_ID_1 (83)<br/>Basic<br/>"]
    q84["Interface__Basic_0 (84)<br/>Basic<br/>"]
    q85{"Interface__LoopEntry_0 (85)<br/>LoopEntry<br/><br/>dec=3"}
    q86["Interface__LoopEnd_0 (86)<br/>LoopEnd<br/>"]
    q87["Interface__LoopBack_0 (87)<br/>LoopBack<br/>"]
    q88{"Interface__Basic_1 (88)<br/>Basic<br/><br/>dec=4"}
    q89["Interface_LeftBrace (89)<br/>Basic<br/>"]
    q90["Interface__Basic_2 (90)<br/>Basic<br/>"]
    q91["Interface__Basic_3 (91)<br/>Basic<br/>"]
    q92{"Interface__LoopEntry_1 (92)<br/>LoopEntry<br/><br/>dec=5"}
    q93["Interface__LoopEnd_1 (93)<br/>LoopEnd<br/>"]
    q94["Interface__LoopBack_1 (94)<br/>LoopBack<br/>"]
    q95["Interface_RightBrace (95)<br/>Basic<br/>"]
    q96["Interface__Basic_4 (96)<br/>Basic<br/>"]

    q2 --> q78
    q78 -->|"tok(&quot;interface&quot;)"| q79
    q79 -->|"tok(ID)"| q88
    q80 -->|"tok(&quot;extends&quot;)"| q81
    q81 -->|"tok(ID)"| q85
    q82 -->|"tok(&quot;,&quot;)"| q83
    q83 -->|"tok(ID)"| q84
    q84 --> q87
    q85 --> q82
    q85 --> q86
    q86 --> q89
    q87 --> q85
    q88 --> q80
    q88 --> q86
    q89 -->|"tok(&quot;{&quot;)"| q92
    q90 -.->|"[Field]"| q91
    q91 --> q94
    q92 --> q90
    q92 --> q93
    q93 --> q95
    q94 --> q92
    q95 -->|"tok(&quot;}&quot;)"| q96
    q96 --> q3
```

## Field

```mermaid
flowchart TD
    q4(["Field__Start (4)<br/>RuleStart"])
    q5(["Field__Stop (5)<br/>RuleStop"])
    q97["Field_Name_ID (97)<br/>Basic<br/>"]
    q98["Field__Basic_0 (98)<br/>Basic<br/>"]
    q99["Field__Basic_1 (99)<br/>Basic<br/>"]

    q4 --> q97
    q97 -->|"tok(ID)"| q98
    q98 -.->|"[FieldType]"| q99
    q99 --> q5
```

## FieldType

```mermaid
flowchart TD
    q6(["FieldType__Start (6)<br/>RuleStart"])
    q7(["FieldType__Stop (7)<br/>RuleStop"])
    q100["FieldType__Basic_0 (100)<br/>Basic<br/>"]
    q101["FieldType__Basic_1 (101)<br/>Basic<br/>"]
    q102["FieldType__Basic_2 (102)<br/>Basic<br/>"]
    q103["FieldType__Basic_3 (103)<br/>Basic<br/>"]
    q104["FieldType__Basic_4 (104)<br/>Basic<br/>"]
    q105["FieldType__Basic_5 (105)<br/>Basic<br/>"]
    q106["FieldType__Basic_6 (106)<br/>Basic<br/>"]
    q107["FieldType__Basic_7 (107)<br/>Basic<br/>"]
    q108{"FieldType__Basic_8 (108)<br/>Basic<br/><br/>dec=6"}
    q109["FieldType__BlockEnd (109)<br/>BlockEnd<br/>"]

    q6 --> q108
    q100 -.->|"[SimpleType]"| q101
    q101 --> q109
    q102 -.->|"[ReferenceType]"| q103
    q103 --> q109
    q104 -.->|"[ArrayType]"| q105
    q105 --> q109
    q106 -.->|"[PrimitiveType]"| q107
    q107 --> q109
    q108 --> q100
    q108 --> q102
    q108 --> q104
    q108 --> q106
    q109 --> q7
```

## ArrayType

```mermaid
flowchart TD
    q8(["ArrayType__Start (8)<br/>RuleStart"])
    q9(["ArrayType__Stop (9)<br/>RuleStop"])
    q110["ArrayType_LeftBracket (110)<br/>Basic<br/>"]
    q111["ArrayType_RightBracket (111)<br/>Basic<br/>"]
    q112["ArrayType__Basic_0 (112)<br/>Basic<br/>"]
    q113["ArrayType__Basic_1 (113)<br/>Basic<br/>"]

    q8 --> q110
    q110 -->|"tok(&quot;[&quot;)"| q111
    q111 -->|"tok(&quot;]&quot;)"| q112
    q112 -.->|"[FieldType]"| q113
    q113 --> q9
```

## ReferenceType

```mermaid
flowchart TD
    q10(["ReferenceType__Start (10)<br/>RuleStart"])
    q11(["ReferenceType__Stop (11)<br/>RuleStop"])
    q114["ReferenceType_Asterisk (114)<br/>Basic<br/>"]
    q115["ReferenceType_Type_ID (115)<br/>Basic<br/>"]
    q116["ReferenceType__Basic (116)<br/>Basic<br/>"]

    q10 --> q114
    q114 -->|"tok(&quot;*&quot;)"| q115
    q115 -->|"tok(ID)"| q116
    q116 --> q11
```

## SimpleType

```mermaid
flowchart TD
    q12(["SimpleType__Start (12)<br/>RuleStart"])
    q13(["SimpleType__Stop (13)<br/>RuleStop"])
    q117["SimpleType_Type_ID (117)<br/>Basic<br/>"]
    q118["SimpleType__Basic (118)<br/>Basic<br/>"]

    q12 --> q117
    q117 -->|"tok(ID)"| q118
    q118 --> q13
```

## PrimitiveType

```mermaid
flowchart TD
    q14(["PrimitiveType__Start (14)<br/>RuleStart"])
    q15(["PrimitiveType__Stop (15)<br/>RuleStop"])
    q119["PrimitiveType_Type_string (119)<br/>Basic<br/>"]
    q120["PrimitiveType__Basic_0 (120)<br/>Basic<br/>"]
    q121["PrimitiveType_Type_bool (121)<br/>Basic<br/>"]
    q122["PrimitiveType__Basic_1 (122)<br/>Basic<br/>"]
    q123["PrimitiveType_Type_composite (123)<br/>Basic<br/>"]
    q124["PrimitiveType__Basic_2 (124)<br/>Basic<br/>"]
    q125{"PrimitiveType__Basic_3 (125)<br/>Basic<br/><br/>dec=7"}
    q126["PrimitiveType__BlockEnd (126)<br/>BlockEnd<br/>"]

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
    q16(["ParserRule__Start (16)<br/>RuleStart"])
    q17(["ParserRule__Stop (17)<br/>RuleStop"])
    q127["ParserRule_Entry_entry (127)<br/>Basic<br/>"]
    q128["ParserRule__Basic_0 (128)<br/>Basic<br/>"]
    q129{"ParserRule__Basic_1 (129)<br/>Basic<br/><br/>dec=8"}
    q130["ParserRule_Name_ID (130)<br/>Basic<br/>"]
    q131["ParserRule_returns (131)<br/>Basic<br/>"]
    q132["ParserRule_ReturnType_ID (132)<br/>Basic<br/>"]
    q133["ParserRule__Basic_2 (133)<br/>Basic<br/>"]
    q134{"ParserRule__Basic_3 (134)<br/>Basic<br/><br/>dec=9"}
    q135["ParserRule_Colon (135)<br/>Basic<br/>"]
    q136["ParserRule__Basic_4 (136)<br/>Basic<br/>"]
    q137["ParserRule_Semicolon (137)<br/>Basic<br/>"]
    q138["ParserRule__Basic_5 (138)<br/>Basic<br/>"]
    q139{"ParserRule__Basic_6 (139)<br/>Basic<br/><br/>dec=10"}

    q16 --> q129
    q127 -->|"tok(&quot;entry&quot;)"| q128
    q128 --> q130
    q129 --> q127
    q129 --> q128
    q130 -->|"tok(ID)"| q134
    q131 -->|"tok(&quot;returns&quot;)"| q132
    q132 -->|"tok(ID)"| q133
    q133 --> q135
    q134 --> q131
    q134 --> q133
    q135 -->|"tok(&quot;:&quot;)"| q136
    q136 -.->|"[Alternatives]"| q139
    q137 -->|"tok(&quot;;&quot;)"| q138
    q138 --> q17
    q139 --> q137
    q139 --> q138
```

## Token

```mermaid
flowchart TD
    q18(["Token__Start (18)<br/>RuleStart"])
    q19(["Token__Stop (19)<br/>RuleStop"])
    q140["Token_Type_hidden (140)<br/>Basic<br/>"]
    q141["Token__Basic_0 (141)<br/>Basic<br/>"]
    q142["Token_Type_comment (142)<br/>Basic<br/>"]
    q143["Token__Basic_1 (143)<br/>Basic<br/>"]
    q144{"Token__Basic_2 (144)<br/>Basic<br/><br/>dec=11"}
    q145["Token__BlockEnd (145)<br/>BlockEnd<br/>"]
    q146{"Token__Basic_3 (146)<br/>Basic<br/><br/>dec=12"}
    q147["Token_token (147)<br/>Basic<br/>"]
    q148["Token_Name_ID (148)<br/>Basic<br/>"]
    q149["Token_Colon (149)<br/>Basic<br/>"]
    q150["Token_Regexp_RegexLiteral (150)<br/>Basic<br/>"]
    q151["Token_DashGreaterThan (151)<br/>Basic<br/>"]
    q152["Token__Basic_4 (152)<br/>Basic<br/>"]
    q153["Token__Basic_5 (153)<br/>Basic<br/>"]
    q154{"Token__Basic_6 (154)<br/>Basic<br/><br/>dec=13"}
    q155["Token_Semicolon (155)<br/>Basic<br/>"]
    q156["Token__Basic_7 (156)<br/>Basic<br/>"]
    q157{"Token__Basic_8 (157)<br/>Basic<br/><br/>dec=14"}

    q18 --> q146
    q140 -->|"tok(&quot;hidden&quot;)"| q141
    q141 --> q145
    q142 -->|"tok(&quot;comment&quot;)"| q143
    q143 --> q145
    q144 --> q140
    q144 --> q142
    q145 --> q147
    q146 --> q144
    q146 --> q145
    q147 -->|"tok(&quot;token&quot;)"| q148
    q148 -->|"tok(ID)"| q149
    q149 -->|"tok(&quot;:&quot;)"| q150
    q150 -->|"tok(RegexLiteral)"| q154
    q151 -->|"tok(&quot;->&quot;)"| q152
    q152 -.->|"[TokenCommand]"| q153
    q153 --> q157
    q154 --> q151
    q154 --> q153
    q155 -->|"tok(&quot;;&quot;)"| q156
    q156 --> q19
    q157 --> q155
    q157 --> q156
```

## TokenCommand

```mermaid
flowchart TD
    q20(["TokenCommand__Start (20)<br/>RuleStart"])
    q21(["TokenCommand__Stop (21)<br/>RuleStop"])
    q158["TokenCommand_Type_push (158)<br/>Basic<br/>"]
    q159["TokenCommand__Basic_0 (159)<br/>Basic<br/>"]
    q160["TokenCommand_Type_pop (160)<br/>Basic<br/>"]
    q161["TokenCommand__Basic_1 (161)<br/>Basic<br/>"]
    q162["TokenCommand_Type_mode (162)<br/>Basic<br/>"]
    q163["TokenCommand__Basic_2 (163)<br/>Basic<br/>"]
    q164{"TokenCommand__Basic_3 (164)<br/>Basic<br/><br/>dec=15"}
    q165["TokenCommand__BlockEnd_0 (165)<br/>BlockEnd<br/>"]
    q166["TokenCommand_LeftParen (166)<br/>Basic<br/>"]
    q167["TokenCommand_Mode_ID (167)<br/>Basic<br/>"]
    q168["TokenCommand__Basic_4 (168)<br/>Basic<br/>"]
    q169["TokenCommand_Default_default (169)<br/>Basic<br/>"]
    q170["TokenCommand__Basic_5 (170)<br/>Basic<br/>"]
    q171{"TokenCommand__Basic_6 (171)<br/>Basic<br/><br/>dec=16"}
    q172["TokenCommand__BlockEnd_1 (172)<br/>BlockEnd<br/>"]
    q173["TokenCommand_RightParen (173)<br/>Basic<br/>"]
    q174["TokenCommand__Basic_7 (174)<br/>Basic<br/>"]
    q175{"TokenCommand__Basic_8 (175)<br/>Basic<br/><br/>dec=17"}

    q20 --> q164
    q158 -->|"tok(&quot;push&quot;)"| q159
    q159 --> q165
    q160 -->|"tok(&quot;pop&quot;)"| q161
    q161 --> q165
    q162 -->|"tok(&quot;mode&quot;)"| q163
    q163 --> q165
    q164 --> q158
    q164 --> q160
    q164 --> q162
    q165 --> q175
    q166 -->|"tok(&quot;(&quot;)"| q171
    q167 -->|"tok(ID)"| q168
    q168 --> q172
    q169 -->|"tok(&quot;default&quot;)"| q170
    q170 --> q172
    q171 --> q167
    q171 --> q169
    q172 --> q173
    q173 -->|"tok(&quot;)&quot;)"| q174
    q174 --> q21
    q175 --> q166
    q175 --> q174
```

## TokenGroup

```mermaid
flowchart TD
    q22(["TokenGroup__Start (22)<br/>RuleStart"])
    q23(["TokenGroup__Stop (23)<br/>RuleStop"])
    q176["TokenGroup_token (176)<br/>Basic<br/>"]
    q177["TokenGroup_group (177)<br/>Basic<br/>"]
    q178["TokenGroup_Name_ID (178)<br/>Basic<br/>"]
    q179["TokenGroup_LeftBrace (179)<br/>Basic<br/>"]
    q180["TokenGroup_TokenRefs_ID (180)<br/>Basic<br/>"]
    q181["TokenGroup__Basic_0 (181)<br/>Basic<br/>"]
    q182["TokenGroup_keywords (182)<br/>Basic<br/>"]
    q183["TokenGroup_Regexps_RegexLiteral (183)<br/>Basic<br/>"]
    q184["TokenGroup__Basic_1 (184)<br/>Basic<br/>"]
    q185["TokenGroup__Basic_2 (185)<br/>Basic<br/>"]
    q186["TokenGroup__Basic_3 (186)<br/>Basic<br/>"]
    q187{"TokenGroup__Basic_4 (187)<br/>Basic<br/><br/>dec=18"}
    q188["TokenGroup__BlockEnd (188)<br/>BlockEnd<br/>"]
    q189{"TokenGroup__LoopEntry (189)<br/>LoopEntry<br/><br/>dec=19"}
    q190["TokenGroup__LoopEnd (190)<br/>LoopEnd<br/>"]
    q191["TokenGroup__LoopBack (191)<br/>LoopBack<br/>"]
    q192["TokenGroup_RightBrace (192)<br/>Basic<br/>"]
    q193["TokenGroup__Basic_5 (193)<br/>Basic<br/>"]

    q22 --> q176
    q176 -->|"tok(&quot;token&quot;)"| q177
    q177 -->|"tok(&quot;group&quot;)"| q178
    q178 -->|"tok(ID)"| q179
    q179 -->|"tok(&quot;{&quot;)"| q189
    q180 -->|"tok(ID)"| q181
    q181 --> q188
    q182 -->|"tok(&quot;keywords&quot;)"| q183
    q183 -->|"tok(RegexLiteral)"| q184
    q184 --> q188
    q185 -.->|"[Keyword]"| q186
    q186 --> q188
    q187 --> q180
    q187 --> q182
    q187 --> q185
    q188 --> q191
    q189 --> q187
    q189 --> q190
    q190 --> q192
    q191 --> q189
    q192 -->|"tok(&quot;}&quot;)"| q193
    q193 --> q23
```

## TokenMode

```mermaid
flowchart TD
    q24(["TokenMode__Start (24)<br/>RuleStart"])
    q25(["TokenMode__Stop (25)<br/>RuleStop"])
    q194["TokenMode_token (194)<br/>Basic<br/>"]
    q195["TokenMode_mode (195)<br/>Basic<br/>"]
    q196["TokenMode_Name_ID (196)<br/>Basic<br/>"]
    q197["TokenMode__Basic_0 (197)<br/>Basic<br/>"]
    q198["TokenMode_Default_default (198)<br/>Basic<br/>"]
    q199["TokenMode__Basic_1 (199)<br/>Basic<br/>"]
    q200{"TokenMode__Basic_2 (200)<br/>Basic<br/><br/>dec=20"}
    q201["TokenMode__BlockEnd_0 (201)<br/>BlockEnd<br/>"]
    q202["TokenMode_LeftBrace (202)<br/>Basic<br/>"]
    q203["TokenMode_TokenRefs_ID (203)<br/>Basic<br/>"]
    q204["TokenMode__Basic_3 (204)<br/>Basic<br/>"]
    q205["TokenMode_keywords (205)<br/>Basic<br/>"]
    q206["TokenMode_Regexps_RegexLiteral (206)<br/>Basic<br/>"]
    q207["TokenMode__Basic_4 (207)<br/>Basic<br/>"]
    q208["TokenMode__Basic_5 (208)<br/>Basic<br/>"]
    q209["TokenMode__Basic_6 (209)<br/>Basic<br/>"]
    q210{"TokenMode__Basic_7 (210)<br/>Basic<br/><br/>dec=21"}
    q211["TokenMode__BlockEnd_1 (211)<br/>BlockEnd<br/>"]
    q212{"TokenMode__LoopEntry (212)<br/>LoopEntry<br/><br/>dec=22"}
    q213["TokenMode__LoopEnd (213)<br/>LoopEnd<br/>"]
    q214["TokenMode__LoopBack (214)<br/>LoopBack<br/>"]
    q215["TokenMode_RightBrace (215)<br/>Basic<br/>"]
    q216["TokenMode__Basic_8 (216)<br/>Basic<br/>"]

    q24 --> q194
    q194 -->|"tok(&quot;token&quot;)"| q195
    q195 -->|"tok(&quot;mode&quot;)"| q200
    q196 -->|"tok(ID)"| q197
    q197 --> q201
    q198 -->|"tok(&quot;default&quot;)"| q199
    q199 --> q201
    q200 --> q196
    q200 --> q198
    q201 --> q202
    q202 -->|"tok(&quot;{&quot;)"| q212
    q203 -->|"tok(ID)"| q204
    q204 --> q211
    q205 -->|"tok(&quot;keywords&quot;)"| q206
    q206 -->|"tok(RegexLiteral)"| q207
    q207 --> q211
    q208 -.->|"[Keyword]"| q209
    q209 --> q211
    q210 --> q203
    q210 --> q205
    q210 --> q208
    q211 --> q214
    q212 --> q210
    q212 --> q213
    q213 --> q215
    q214 --> q212
    q215 -->|"tok(&quot;}&quot;)"| q216
    q216 --> q25
```

## Alternatives

```mermaid
flowchart TD
    q26(["Alternatives__Start (26)<br/>RuleStart"])
    q27(["Alternatives__Stop (27)<br/>RuleStop"])
    q217["Alternatives__Basic_0 (217)<br/>Basic<br/>"]
    q218["Alternatives_Pipe (218)<br/>Basic<br/>"]
    q219["Alternatives__Basic_1 (219)<br/>Basic<br/>"]
    q220["Alternatives__Basic_2 (220)<br/>Basic<br/>"]
    q221{"Alternatives__LoopBack (221)<br/>LoopBack<br/><br/>dec=23"}
    q222["Alternatives__LoopEnd (222)<br/>LoopEnd<br/>"]
    q223{"Alternatives__Basic_3 (223)<br/>Basic<br/><br/>dec=24"}

    q26 --> q217
    q217 -.->|"[Group]"| q223
    q218 -->|"tok(&quot;|&quot;)"| q219
    q219 -.->|"[Group]"| q220
    q220 --> q221
    q221 --> q218
    q221 --> q222
    q222 --> q27
    q223 --> q218
    q223 --> q222
```

## Group

```mermaid
flowchart TD
    q28(["Group__Start (28)<br/>RuleStart"])
    q29(["Group__Stop (29)<br/>RuleStop"])
    q224["Group__Basic_0 (224)<br/>Basic<br/>"]
    q225["Group__Basic_1 (225)<br/>Basic<br/>"]
    q226["Group__Basic_2 (226)<br/>Basic<br/>"]
    q227{"Group__LoopBack (227)<br/>LoopBack<br/><br/>dec=25"}
    q228["Group__LoopEnd (228)<br/>LoopEnd<br/>"]
    q229{"Group__Basic_3 (229)<br/>Basic<br/><br/>dec=26"}

    q28 --> q224
    q224 -.->|"[Element]"| q229
    q225 -.->|"[Element]"| q226
    q226 --> q227
    q227 --> q225
    q227 --> q228
    q228 --> q29
    q229 --> q225
    q229 --> q228
```

## Element

```mermaid
flowchart TD
    q30(["Element__Start (30)<br/>RuleStart"])
    q31(["Element__Stop (31)<br/>RuleStop"])
    q230["Element__Basic_0 (230)<br/>Basic<br/>"]
    q231["Element__Basic_1 (231)<br/>Basic<br/>"]
    q232["Element__Basic_2 (232)<br/>Basic<br/>"]
    q233["Element__Basic_3 (233)<br/>Basic<br/>"]
    q234["Element__Basic_4 (234)<br/>Basic<br/>"]
    q235["Element__Basic_5 (235)<br/>Basic<br/>"]
    q236["Element__Basic_6 (236)<br/>Basic<br/>"]
    q237["Element__Basic_7 (237)<br/>Basic<br/>"]
    q238["Element_LeftParen (238)<br/>Basic<br/>"]
    q239["Element__Basic_8 (239)<br/>Basic<br/>"]
    q240["Element_RightParen (240)<br/>Basic<br/>"]
    q241["Element__Basic_9 (241)<br/>Basic<br/>"]
    q242{"Element__Basic_10 (242)<br/>Basic<br/><br/>dec=27"}
    q243["Element__BlockEnd_0 (243)<br/>BlockEnd<br/>"]
    q244["Element_Cardinality_Asterisk (244)<br/>Basic<br/>"]
    q245["Element__Basic_11 (245)<br/>Basic<br/>"]
    q246["Element_Cardinality_Plus (246)<br/>Basic<br/>"]
    q247["Element__Basic_12 (247)<br/>Basic<br/>"]
    q248["Element_Cardinality_Question (248)<br/>Basic<br/>"]
    q249["Element__Basic_13 (249)<br/>Basic<br/>"]
    q250{"Element__Basic_14 (250)<br/>Basic<br/><br/>dec=28"}
    q251["Element__BlockEnd_1 (251)<br/>BlockEnd<br/>"]
    q252{"Element__Basic_15 (252)<br/>Basic<br/><br/>dec=29"}

    q30 --> q242
    q230 -.->|"[Keyword]"| q231
    q231 --> q243
    q232 -.->|"[Assignment]"| q233
    q233 --> q243
    q234 -.->|"[RuleCall]"| q235
    q235 --> q243
    q236 -.->|"[Action]"| q237
    q237 --> q243
    q238 -->|"tok(&quot;(&quot;)"| q239
    q239 -.->|"[Alternatives]"| q240
    q240 -->|"tok(&quot;)&quot;)"| q241
    q241 --> q243
    q242 --> q230
    q242 --> q232
    q242 --> q234
    q242 --> q236
    q242 --> q238
    q243 --> q252
    q244 -->|"tok(&quot;*&quot;)"| q245
    q245 --> q251
    q246 -->|"tok(&quot;+&quot;)"| q247
    q247 --> q251
    q248 -->|"tok(&quot;?&quot;)"| q249
    q249 --> q251
    q250 --> q244
    q250 --> q246
    q250 --> q248
    q251 --> q31
    q252 --> q250
    q252 --> q251
```

## Keyword

```mermaid
flowchart TD
    q32(["Keyword__Start (32)<br/>RuleStart"])
    q33(["Keyword__Stop (33)<br/>RuleStop"])
    q253["Keyword_Value_StringLiteral (253)<br/>Basic<br/>"]
    q254["Keyword__Basic (254)<br/>Basic<br/>"]

    q32 --> q253
    q253 -->|"tok(StringLiteral)"| q254
    q254 --> q33
```

## Assignment

```mermaid
flowchart TD
    q34(["Assignment__Start (34)<br/>RuleStart"])
    q35(["Assignment__Stop (35)<br/>RuleStop"])
    q255["Assignment_Property_ID (255)<br/>Basic<br/>"]
    q256["Assignment_Operator_PlusEquals (256)<br/>Basic<br/>"]
    q257["Assignment__Basic_0 (257)<br/>Basic<br/>"]
    q258["Assignment_Operator_Equals (258)<br/>Basic<br/>"]
    q259["Assignment__Basic_1 (259)<br/>Basic<br/>"]
    q260["Assignment_Operator_QuestionEquals (260)<br/>Basic<br/>"]
    q261["Assignment__Basic_2 (261)<br/>Basic<br/>"]
    q262{"Assignment__Basic_3 (262)<br/>Basic<br/><br/>dec=30"}
    q263["Assignment__BlockEnd (263)<br/>BlockEnd<br/>"]
    q264["Assignment__Basic_4 (264)<br/>Basic<br/>"]
    q265["Assignment__Basic_5 (265)<br/>Basic<br/>"]

    q34 --> q255
    q255 -->|"tok(ID)"| q262
    q256 -->|"tok(&quot;+=&quot;)"| q257
    q257 --> q263
    q258 -->|"tok(&quot;=&quot;)"| q259
    q259 --> q263
    q260 -->|"tok(&quot;?=&quot;)"| q261
    q261 --> q263
    q262 --> q256
    q262 --> q258
    q262 --> q260
    q263 --> q264
    q264 -.->|"[Assignable]"| q265
    q265 --> q35
```

## Assignable

```mermaid
flowchart TD
    q36(["Assignable__Start (36)<br/>RuleStart"])
    q37(["Assignable__Stop (37)<br/>RuleStop"])
    q266["Assignable__Basic_0 (266)<br/>Basic<br/>"]
    q267["Assignable__Basic_1 (267)<br/>Basic<br/>"]
    q268["Assignable__Basic_2 (268)<br/>Basic<br/>"]
    q269["Assignable__Basic_3 (269)<br/>Basic<br/>"]
    q270["Assignable__Basic_4 (270)<br/>Basic<br/>"]
    q271["Assignable__Basic_5 (271)<br/>Basic<br/>"]
    q272["Assignable_LeftParen (272)<br/>Basic<br/>"]
    q273["Assignable__Basic_6 (273)<br/>Basic<br/>"]
    q274["Assignable_RightParen (274)<br/>Basic<br/>"]
    q275["Assignable__Basic_7 (275)<br/>Basic<br/>"]
    q276{"Assignable__Basic_8 (276)<br/>Basic<br/><br/>dec=31"}
    q277["Assignable__BlockEnd (277)<br/>BlockEnd<br/>"]

    q36 --> q276
    q266 -.->|"[Keyword]"| q267
    q267 --> q277
    q268 -.->|"[RuleCall]"| q269
    q269 --> q277
    q270 -.->|"[CrossRef]"| q271
    q271 --> q277
    q272 -->|"tok(&quot;(&quot;)"| q273
    q273 -.->|"[AssignableAlternatives]"| q274
    q274 -->|"tok(&quot;)&quot;)"| q275
    q275 --> q277
    q276 --> q266
    q276 --> q268
    q276 --> q270
    q276 --> q272
    q277 --> q37
```

## AssignableWithoutAlts

```mermaid
flowchart TD
    q38(["AssignableWithoutAlts__Start (38)<br/>RuleStart"])
    q39(["AssignableWithoutAlts__Stop (39)<br/>RuleStop"])
    q278["AssignableWithoutAlts__Basic_0 (278)<br/>Basic<br/>"]
    q279["AssignableWithoutAlts__Basic_1 (279)<br/>Basic<br/>"]
    q280["AssignableWithoutAlts__Basic_2 (280)<br/>Basic<br/>"]
    q281["AssignableWithoutAlts__Basic_3 (281)<br/>Basic<br/>"]
    q282["AssignableWithoutAlts__Basic_4 (282)<br/>Basic<br/>"]
    q283["AssignableWithoutAlts__Basic_5 (283)<br/>Basic<br/>"]
    q284{"AssignableWithoutAlts__Basic_6 (284)<br/>Basic<br/><br/>dec=32"}
    q285["AssignableWithoutAlts__BlockEnd (285)<br/>BlockEnd<br/>"]

    q38 --> q284
    q278 -.->|"[Keyword]"| q279
    q279 --> q285
    q280 -.->|"[RuleCall]"| q281
    q281 --> q285
    q282 -.->|"[CrossRef]"| q283
    q283 --> q285
    q284 --> q278
    q284 --> q280
    q284 --> q282
    q285 --> q39
```

## AssignableAlternatives

```mermaid
flowchart TD
    q40(["AssignableAlternatives__Start (40)<br/>RuleStart"])
    q41(["AssignableAlternatives__Stop (41)<br/>RuleStop"])
    q286["AssignableAlternatives__Basic_0 (286)<br/>Basic<br/>"]
    q287["AssignableAlternatives_Pipe (287)<br/>Basic<br/>"]
    q288["AssignableAlternatives__Basic_1 (288)<br/>Basic<br/>"]
    q289["AssignableAlternatives__Basic_2 (289)<br/>Basic<br/>"]
    q290{"AssignableAlternatives__LoopBack (290)<br/>LoopBack<br/><br/>dec=33"}
    q291["AssignableAlternatives__LoopEnd (291)<br/>LoopEnd<br/>"]
    q292{"AssignableAlternatives__Basic_3 (292)<br/>Basic<br/><br/>dec=34"}

    q40 --> q286
    q286 -.->|"[AssignableWithoutAlts]"| q292
    q287 -->|"tok(&quot;|&quot;)"| q288
    q288 -.->|"[AssignableWithoutAlts]"| q289
    q289 --> q290
    q290 --> q287
    q290 --> q291
    q291 --> q41
    q292 --> q287
    q292 --> q291
```

## CrossRef

```mermaid
flowchart TD
    q42(["CrossRef__Start (42)<br/>RuleStart"])
    q43(["CrossRef__Stop (43)<br/>RuleStop"])
    q293["CrossRef_LeftBracket (293)<br/>Basic<br/>"]
    q294["CrossRef_Type_ID (294)<br/>Basic<br/>"]
    q295["CrossRef_Colon (295)<br/>Basic<br/>"]
    q296["CrossRef__Basic_0 (296)<br/>Basic<br/>"]
    q297["CrossRef__Basic_1 (297)<br/>Basic<br/>"]
    q298{"CrossRef__Basic_2 (298)<br/>Basic<br/><br/>dec=35"}
    q299["CrossRef_RightBracket (299)<br/>Basic<br/>"]
    q300["CrossRef__Basic_3 (300)<br/>Basic<br/>"]

    q42 --> q293
    q293 -->|"tok(&quot;[&quot;)"| q294
    q294 -->|"tok(ID)"| q298
    q295 -->|"tok(&quot;:&quot;)"| q296
    q296 -.->|"[RuleCall]"| q297
    q297 --> q299
    q298 --> q295
    q298 --> q297
    q299 -->|"tok(&quot;]&quot;)"| q300
    q300 --> q43
```

## RuleCall

```mermaid
flowchart TD
    q44(["RuleCall__Start (44)<br/>RuleStart"])
    q45(["RuleCall__Stop (45)<br/>RuleStop"])
    q301["RuleCall_Rule_ID (301)<br/>Basic<br/>"]
    q302["RuleCall__Basic (302)<br/>Basic<br/>"]

    q44 --> q301
    q301 -->|"tok(ID)"| q302
    q302 --> q45
```

## Action

```mermaid
flowchart TD
    q46(["Action__Start (46)<br/>RuleStart"])
    q47(["Action__Stop (47)<br/>RuleStop"])
    q303["Action_LeftBrace (303)<br/>Basic<br/>"]
    q304["Action_Type_ID (304)<br/>Basic<br/>"]
    q305["Action_Dot (305)<br/>Basic<br/>"]
    q306["Action_Property_ID (306)<br/>Basic<br/>"]
    q307["Action_Operator_PlusEquals (307)<br/>Basic<br/>"]
    q308["Action__Basic_0 (308)<br/>Basic<br/>"]
    q309["Action_Operator_Equals (309)<br/>Basic<br/>"]
    q310["Action__Basic_1 (310)<br/>Basic<br/>"]
    q311{"Action__Basic_2 (311)<br/>Basic<br/><br/>dec=36"}
    q312["Action__BlockEnd (312)<br/>BlockEnd<br/>"]
    q313["Action_current (313)<br/>Basic<br/>"]
    q314["Action__Basic_3 (314)<br/>Basic<br/>"]
    q315{"Action__Basic_4 (315)<br/>Basic<br/><br/>dec=37"}
    q316["Action_RightBrace (316)<br/>Basic<br/>"]
    q317["Action__Basic_5 (317)<br/>Basic<br/>"]

    q46 --> q303
    q303 -->|"tok(&quot;{&quot;)"| q304
    q304 -->|"tok(ID)"| q315
    q305 -->|"tok(&quot;.&quot;)"| q306
    q306 -->|"tok(ID)"| q311
    q307 -->|"tok(&quot;+=&quot;)"| q308
    q308 --> q312
    q309 -->|"tok(&quot;=&quot;)"| q310
    q310 --> q312
    q311 --> q307
    q311 --> q309
    q312 --> q313
    q313 -->|"tok(&quot;current&quot;)"| q314
    q314 --> q316
    q315 --> q305
    q315 --> q314
    q316 -->|"tok(&quot;}&quot;)"| q317
    q317 --> q47
```

## CompositeRule

```mermaid
flowchart TD
    q48(["CompositeRule__Start (48)<br/>RuleStart"])
    q49(["CompositeRule__Stop (49)<br/>RuleStop"])
    q318["CompositeRule_composite (318)<br/>Basic<br/>"]
    q319["CompositeRule_Name_ID (319)<br/>Basic<br/>"]
    q320["CompositeRule_Colon (320)<br/>Basic<br/>"]
    q321["CompositeRule__Basic_0 (321)<br/>Basic<br/>"]
    q322["CompositeRule_Semicolon (322)<br/>Basic<br/>"]
    q323["CompositeRule__Basic_1 (323)<br/>Basic<br/>"]
    q324{"CompositeRule__Basic_2 (324)<br/>Basic<br/><br/>dec=38"}

    q48 --> q318
    q318 -->|"tok(&quot;composite&quot;)"| q319
    q319 -->|"tok(ID)"| q320
    q320 -->|"tok(&quot;:&quot;)"| q321
    q321 -.->|"[CompositeAlternatives]"| q324
    q322 -->|"tok(&quot;;&quot;)"| q323
    q323 --> q49
    q324 --> q322
    q324 --> q323
```

## CompositeAlternatives

```mermaid
flowchart TD
    q50(["CompositeAlternatives__Start (50)<br/>RuleStart"])
    q51(["CompositeAlternatives__Stop (51)<br/>RuleStop"])
    q325["CompositeAlternatives__Basic_0 (325)<br/>Basic<br/>"]
    q326["CompositeAlternatives_Pipe (326)<br/>Basic<br/>"]
    q327["CompositeAlternatives__Basic_1 (327)<br/>Basic<br/>"]
    q328["CompositeAlternatives__Basic_2 (328)<br/>Basic<br/>"]
    q329{"CompositeAlternatives__LoopBack (329)<br/>LoopBack<br/><br/>dec=39"}
    q330["CompositeAlternatives__LoopEnd (330)<br/>LoopEnd<br/>"]
    q331{"CompositeAlternatives__Basic_3 (331)<br/>Basic<br/><br/>dec=40"}

    q50 --> q325
    q325 -.->|"[CompositeGroup]"| q331
    q326 -->|"tok(&quot;|&quot;)"| q327
    q327 -.->|"[CompositeGroup]"| q328
    q328 --> q329
    q329 --> q326
    q329 --> q330
    q330 --> q51
    q331 --> q326
    q331 --> q330
```

## CompositeGroup

```mermaid
flowchart TD
    q52(["CompositeGroup__Start (52)<br/>RuleStart"])
    q53(["CompositeGroup__Stop (53)<br/>RuleStop"])
    q332["CompositeGroup__Basic_0 (332)<br/>Basic<br/>"]
    q333["CompositeGroup__Basic_1 (333)<br/>Basic<br/>"]
    q334["CompositeGroup__Basic_2 (334)<br/>Basic<br/>"]
    q335{"CompositeGroup__LoopBack (335)<br/>LoopBack<br/><br/>dec=41"}
    q336["CompositeGroup__LoopEnd (336)<br/>LoopEnd<br/>"]
    q337{"CompositeGroup__Basic_3 (337)<br/>Basic<br/><br/>dec=42"}

    q52 --> q332
    q332 -.->|"[CompositeElement]"| q337
    q333 -.->|"[CompositeElement]"| q334
    q334 --> q335
    q335 --> q333
    q335 --> q336
    q336 --> q53
    q337 --> q333
    q337 --> q336
```

## CompositeElement

```mermaid
flowchart TD
    q54(["CompositeElement__Start (54)<br/>RuleStart"])
    q55(["CompositeElement__Stop (55)<br/>RuleStop"])
    q338["CompositeElement__Basic_0 (338)<br/>Basic<br/>"]
    q339["CompositeElement__Basic_1 (339)<br/>Basic<br/>"]
    q340["CompositeElement__Basic_2 (340)<br/>Basic<br/>"]
    q341["CompositeElement__Basic_3 (341)<br/>Basic<br/>"]
    q342["CompositeElement_LeftParen (342)<br/>Basic<br/>"]
    q343["CompositeElement__Basic_4 (343)<br/>Basic<br/>"]
    q344["CompositeElement_RightParen (344)<br/>Basic<br/>"]
    q345["CompositeElement__Basic_5 (345)<br/>Basic<br/>"]
    q346{"CompositeElement__Basic_6 (346)<br/>Basic<br/><br/>dec=43"}
    q347["CompositeElement__BlockEnd_0 (347)<br/>BlockEnd<br/>"]
    q348["CompositeElement_Cardinality_Asterisk (348)<br/>Basic<br/>"]
    q349["CompositeElement__Basic_7 (349)<br/>Basic<br/>"]
    q350["CompositeElement_Cardinality_Plus (350)<br/>Basic<br/>"]
    q351["CompositeElement__Basic_8 (351)<br/>Basic<br/>"]
    q352["CompositeElement_Cardinality_Question (352)<br/>Basic<br/>"]
    q353["CompositeElement__Basic_9 (353)<br/>Basic<br/>"]
    q354{"CompositeElement__Basic_10 (354)<br/>Basic<br/><br/>dec=44"}
    q355["CompositeElement__BlockEnd_1 (355)<br/>BlockEnd<br/>"]
    q356{"CompositeElement__Basic_11 (356)<br/>Basic<br/><br/>dec=45"}

    q54 --> q346
    q338 -.->|"[Keyword]"| q339
    q339 --> q347
    q340 -.->|"[RuleCall]"| q341
    q341 --> q347
    q342 -->|"tok(&quot;(&quot;)"| q343
    q343 -.->|"[CompositeAlternatives]"| q344
    q344 -->|"tok(&quot;)&quot;)"| q345
    q345 --> q347
    q346 --> q338
    q346 --> q340
    q346 --> q342
    q347 --> q356
    q348 -->|"tok(&quot;*&quot;)"| q349
    q349 --> q355
    q350 -->|"tok(&quot;+&quot;)"| q351
    q351 --> q355
    q352 -->|"tok(&quot;?&quot;)"| q353
    q353 --> q355
    q354 --> q348
    q354 --> q350
    q354 --> q352
    q355 --> q55
    q356 --> q354
    q356 --> q355
```

