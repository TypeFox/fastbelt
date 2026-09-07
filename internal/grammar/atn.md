# Runtime ATN for grammar

## Grammar

```mermaid
flowchart TD
    q0(["Grammar__Start (0)<br/>RuleStart"])
    q1(["Grammar__Stop (1)<br/>RuleStop"])
    q80["Grammar_grammar (80)<br/>Basic<br/>"]
    q81["Grammar_Name_ID (81)<br/>Basic<br/>"]
    q82["Grammar_Semicolon (82)<br/>Basic<br/>"]
    q83["Grammar__Basic_0 (83)<br/>Basic<br/>"]
    q84{"Grammar__Basic_1 (84)<br/>Basic<br/><br/>dec=0"}
    q85["Grammar__Basic_2 (85)<br/>Basic<br/>"]
    q86["Grammar__Basic_3 (86)<br/>Basic<br/>"]
    q87["Grammar__Basic_4 (87)<br/>Basic<br/>"]
    q88["Grammar__Basic_5 (88)<br/>Basic<br/>"]
    q89["Grammar__Basic_6 (89)<br/>Basic<br/>"]
    q90["Grammar__Basic_7 (90)<br/>Basic<br/>"]
    q91["Grammar__Basic_8 (91)<br/>Basic<br/>"]
    q92["Grammar__Basic_9 (92)<br/>Basic<br/>"]
    q93["Grammar__Basic_10 (93)<br/>Basic<br/>"]
    q94["Grammar__Basic_11 (94)<br/>Basic<br/>"]
    q95["Grammar__Basic_12 (95)<br/>Basic<br/>"]
    q96["Grammar__Basic_13 (96)<br/>Basic<br/>"]
    q97["Grammar__Basic_14 (97)<br/>Basic<br/>"]
    q98["Grammar__Basic_15 (98)<br/>Basic<br/>"]
    q99{"Grammar__Basic_16 (99)<br/>Basic<br/><br/>dec=1"}
    q100["Grammar__BlockEnd (100)<br/>BlockEnd<br/>"]
    q101{"Grammar__LoopEntry (101)<br/>LoopEntry<br/><br/>dec=2"}
    q102["Grammar__LoopEnd (102)<br/>LoopEnd<br/>"]
    q103["Grammar__LoopBack (103)<br/>LoopBack<br/>"]

    q0 --> q80
    q80 -->|"tok(&quot;grammar&quot;)"| q81
    q81 -->|"tok(ID)"| q84
    q82 -->|"tok(&quot;;&quot;)"| q83
    q83 --> q101
    q84 --> q82
    q84 --> q83
    q85 -.->|"[ParserRule]"| q86
    q86 --> q100
    q87 -.->|"[TokenDecl]"| q88
    q88 --> q100
    q89 -.->|"[TokenGroup]"| q90
    q90 --> q100
    q91 -.->|"[TokenMode]"| q92
    q92 --> q100
    q93 -.->|"[Interface]"| q94
    q94 --> q100
    q95 -.->|"[CompositeRule]"| q96
    q96 --> q100
    q97 -.->|"[InfixRule]"| q98
    q98 --> q100
    q99 --> q85
    q99 --> q87
    q99 --> q89
    q99 --> q91
    q99 --> q93
    q99 --> q95
    q99 --> q97
    q100 --> q103
    q101 --> q99
    q101 --> q102
    q102 --> q1
    q103 --> q101
```

## Interface

```mermaid
flowchart TD
    q2(["Interface__Start (2)<br/>RuleStart"])
    q3(["Interface__Stop (3)<br/>RuleStop"])
    q104["Interface_interface (104)<br/>Basic<br/>"]
    q105["Interface_Name_ID (105)<br/>Basic<br/>"]
    q106["Interface_extends (106)<br/>Basic<br/>"]
    q107["Interface_Extends_ID_0 (107)<br/>Basic<br/>"]
    q108["Interface_Comma (108)<br/>Basic<br/>"]
    q109["Interface_Extends_ID_1 (109)<br/>Basic<br/>"]
    q110["Interface__Basic_0 (110)<br/>Basic<br/>"]
    q111{"Interface__LoopEntry_0 (111)<br/>LoopEntry<br/><br/>dec=3"}
    q112["Interface__LoopEnd_0 (112)<br/>LoopEnd<br/>"]
    q113["Interface__LoopBack_0 (113)<br/>LoopBack<br/>"]
    q114{"Interface__Basic_1 (114)<br/>Basic<br/><br/>dec=4"}
    q115["Interface_LeftBrace (115)<br/>Basic<br/>"]
    q116["Interface__Basic_2 (116)<br/>Basic<br/>"]
    q117["Interface__Basic_3 (117)<br/>Basic<br/>"]
    q118{"Interface__LoopEntry_1 (118)<br/>LoopEntry<br/><br/>dec=5"}
    q119["Interface__LoopEnd_1 (119)<br/>LoopEnd<br/>"]
    q120["Interface__LoopBack_1 (120)<br/>LoopBack<br/>"]
    q121["Interface_RightBrace (121)<br/>Basic<br/>"]
    q122["Interface__Basic_4 (122)<br/>Basic<br/>"]

    q2 --> q104
    q104 -->|"tok(&quot;interface&quot;)"| q105
    q105 -->|"tok(ID)"| q114
    q106 -->|"tok(&quot;extends&quot;)"| q107
    q107 -->|"tok(ID)"| q111
    q108 -->|"tok(&quot;,&quot;)"| q109
    q109 -->|"tok(ID)"| q110
    q110 --> q113
    q111 --> q108
    q111 --> q112
    q112 --> q115
    q113 --> q111
    q114 --> q106
    q114 --> q112
    q115 -->|"tok(&quot;{&quot;)"| q118
    q116 -.->|"[Field]"| q117
    q117 --> q120
    q118 --> q116
    q118 --> q119
    q119 --> q121
    q120 --> q118
    q121 -->|"tok(&quot;}&quot;)"| q122
    q122 --> q3
```

## Field

```mermaid
flowchart TD
    q4(["Field__Start (4)<br/>RuleStart"])
    q5(["Field__Stop (5)<br/>RuleStop"])
    q123["Field_Name_ID (123)<br/>Basic<br/>"]
    q124["Field__Basic_0 (124)<br/>Basic<br/>"]
    q125["Field__Basic_1 (125)<br/>Basic<br/>"]

    q4 --> q123
    q123 -->|"tok(ID)"| q124
    q124 -.->|"[FieldType]"| q125
    q125 --> q5
```

## FieldType

```mermaid
flowchart TD
    q6(["FieldType__Start (6)<br/>RuleStart"])
    q7(["FieldType__Stop (7)<br/>RuleStop"])
    q126["FieldType__Basic_0 (126)<br/>Basic<br/>"]
    q127["FieldType__Basic_1 (127)<br/>Basic<br/>"]
    q128["FieldType__Basic_2 (128)<br/>Basic<br/>"]
    q129["FieldType__Basic_3 (129)<br/>Basic<br/>"]
    q130["FieldType__Basic_4 (130)<br/>Basic<br/>"]
    q131["FieldType__Basic_5 (131)<br/>Basic<br/>"]
    q132["FieldType__Basic_6 (132)<br/>Basic<br/>"]
    q133["FieldType__Basic_7 (133)<br/>Basic<br/>"]
    q134{"FieldType__Basic_8 (134)<br/>Basic<br/><br/>dec=6"}
    q135["FieldType__BlockEnd (135)<br/>BlockEnd<br/>"]

    q6 --> q134
    q126 -.->|"[SimpleType]"| q127
    q127 --> q135
    q128 -.->|"[ReferenceType]"| q129
    q129 --> q135
    q130 -.->|"[ArrayType]"| q131
    q131 --> q135
    q132 -.->|"[PrimitiveType]"| q133
    q133 --> q135
    q134 --> q126
    q134 --> q128
    q134 --> q130
    q134 --> q132
    q135 --> q7
```

## ArrayType

```mermaid
flowchart TD
    q8(["ArrayType__Start (8)<br/>RuleStart"])
    q9(["ArrayType__Stop (9)<br/>RuleStop"])
    q136["ArrayType_LeftBracket (136)<br/>Basic<br/>"]
    q137["ArrayType_RightBracket (137)<br/>Basic<br/>"]
    q138["ArrayType__Basic_0 (138)<br/>Basic<br/>"]
    q139["ArrayType__Basic_1 (139)<br/>Basic<br/>"]

    q8 --> q136
    q136 -->|"tok(&quot;[&quot;)"| q137
    q137 -->|"tok(&quot;]&quot;)"| q138
    q138 -.->|"[FieldType]"| q139
    q139 --> q9
```

## ReferenceType

```mermaid
flowchart TD
    q10(["ReferenceType__Start (10)<br/>RuleStart"])
    q11(["ReferenceType__Stop (11)<br/>RuleStop"])
    q140["ReferenceType_Asterisk (140)<br/>Basic<br/>"]
    q141["ReferenceType_Type_ID (141)<br/>Basic<br/>"]
    q142["ReferenceType__Basic (142)<br/>Basic<br/>"]

    q10 --> q140
    q140 -->|"tok(&quot;*&quot;)"| q141
    q141 -->|"tok(ID)"| q142
    q142 --> q11
```

## SimpleType

```mermaid
flowchart TD
    q12(["SimpleType__Start (12)<br/>RuleStart"])
    q13(["SimpleType__Stop (13)<br/>RuleStop"])
    q143["SimpleType_Type_ID (143)<br/>Basic<br/>"]
    q144["SimpleType__Basic (144)<br/>Basic<br/>"]

    q12 --> q143
    q143 -->|"tok(ID)"| q144
    q144 --> q13
```

## PrimitiveType

```mermaid
flowchart TD
    q14(["PrimitiveType__Start (14)<br/>RuleStart"])
    q15(["PrimitiveType__Stop (15)<br/>RuleStop"])
    q145["PrimitiveType_Type_string (145)<br/>Basic<br/>"]
    q146["PrimitiveType__Basic_0 (146)<br/>Basic<br/>"]
    q147["PrimitiveType_Type_bool (147)<br/>Basic<br/>"]
    q148["PrimitiveType__Basic_1 (148)<br/>Basic<br/>"]
    q149["PrimitiveType_Type_composite (149)<br/>Basic<br/>"]
    q150["PrimitiveType__Basic_2 (150)<br/>Basic<br/>"]
    q151{"PrimitiveType__Basic_3 (151)<br/>Basic<br/><br/>dec=7"}
    q152["PrimitiveType__BlockEnd (152)<br/>BlockEnd<br/>"]

    q14 --> q151
    q145 -->|"tok(&quot;string&quot;)"| q146
    q146 --> q152
    q147 -->|"tok(&quot;bool&quot;)"| q148
    q148 --> q152
    q149 -->|"tok(&quot;composite&quot;)"| q150
    q150 --> q152
    q151 --> q145
    q151 --> q147
    q151 --> q149
    q152 --> q15
```

## ParserRule

```mermaid
flowchart TD
    q16(["ParserRule__Start (16)<br/>RuleStart"])
    q17(["ParserRule__Stop (17)<br/>RuleStop"])
    q153["ParserRule_Entry_entry (153)<br/>Basic<br/>"]
    q154["ParserRule__Basic_0 (154)<br/>Basic<br/>"]
    q155{"ParserRule__Basic_1 (155)<br/>Basic<br/><br/>dec=8"}
    q156["ParserRule_Name_ID (156)<br/>Basic<br/>"]
    q157["ParserRule_returns (157)<br/>Basic<br/>"]
    q158["ParserRule_ReturnType_ID (158)<br/>Basic<br/>"]
    q159["ParserRule__Basic_2 (159)<br/>Basic<br/>"]
    q160{"ParserRule__Basic_3 (160)<br/>Basic<br/><br/>dec=9"}
    q161["ParserRule_Colon (161)<br/>Basic<br/>"]
    q162["ParserRule__Basic_4 (162)<br/>Basic<br/>"]
    q163["ParserRule_Semicolon (163)<br/>Basic<br/>"]
    q164["ParserRule__Basic_5 (164)<br/>Basic<br/>"]
    q165{"ParserRule__Basic_6 (165)<br/>Basic<br/><br/>dec=10"}

    q16 --> q155
    q153 -->|"tok(&quot;entry&quot;)"| q154
    q154 --> q156
    q155 --> q153
    q155 --> q154
    q156 -->|"tok(ID)"| q160
    q157 -->|"tok(&quot;returns&quot;)"| q158
    q158 -->|"tok(ID)"| q159
    q159 --> q161
    q160 --> q157
    q160 --> q159
    q161 -->|"tok(&quot;:&quot;)"| q162
    q162 -.->|"[Alternatives]"| q165
    q163 -->|"tok(&quot;;&quot;)"| q164
    q164 --> q17
    q165 --> q163
    q165 --> q164
```

## TokenDecl

```mermaid
flowchart TD
    q18(["TokenDecl__Start (18)<br/>RuleStart"])
    q19(["TokenDecl__Stop (19)<br/>RuleStop"])
    q166["TokenDecl_Modifier_TokenModifier (166)<br/>Basic<br/>"]
    q167["TokenDecl__Basic_0 (167)<br/>Basic<br/>"]
    q168{"TokenDecl__Basic_1 (168)<br/>Basic<br/><br/>dec=11"}
    q169["TokenDecl_token (169)<br/>Basic<br/>"]
    q170["TokenDecl_Name_ID (170)<br/>Basic<br/>"]
    q171["TokenDecl_Colon (171)<br/>Basic<br/>"]
    q172["TokenDecl__Basic_2 (172)<br/>Basic<br/>"]
    q173["TokenDecl__Basic_3 (173)<br/>Basic<br/>"]
    q174["TokenDecl__Basic_4 (174)<br/>Basic<br/>"]
    q175{"TokenDecl__Basic_5 (175)<br/>Basic<br/><br/>dec=12"}
    q176["TokenDecl_Semicolon (176)<br/>Basic<br/>"]
    q177["TokenDecl__Basic_6 (177)<br/>Basic<br/>"]
    q178{"TokenDecl__Basic_7 (178)<br/>Basic<br/><br/>dec=13"}

    q18 --> q168
    q166 -->|"tok(TokenModifier)"| q167
    q167 --> q169
    q168 --> q166
    q168 --> q167
    q169 -->|"tok(&quot;token&quot;)"| q170
    q170 -->|"tok(ID)"| q171
    q171 -->|"tok(&quot;:&quot;)"| q172
    q172 -.->|"[TokenContent]"| q175
    q173 -.->|"[TokenCommand]"| q174
    q174 --> q178
    q175 --> q173
    q175 --> q174
    q176 -->|"tok(&quot;;&quot;)"| q177
    q177 --> q19
    q178 --> q176
    q178 --> q177
```

## TokenContent

```mermaid
flowchart TD
    q20(["TokenContent__Start (20)<br/>RuleStart"])
    q21(["TokenContent__Stop (21)<br/>RuleStop"])
    q179["TokenContent__Basic_0 (179)<br/>Basic<br/>"]
    q180["TokenContent__Basic_1 (180)<br/>Basic<br/>"]
    q181["TokenContent__Basic_2 (181)<br/>Basic<br/>"]
    q182["TokenContent__Basic_3 (182)<br/>Basic<br/>"]
    q183{"TokenContent__Basic_4 (183)<br/>Basic<br/><br/>dec=14"}
    q184["TokenContent__BlockEnd (184)<br/>BlockEnd<br/>"]

    q20 --> q183
    q179 -.->|"[RegexpTokenContent]"| q180
    q180 --> q184
    q181 -.->|"[KeywordTokenContent]"| q182
    q182 --> q184
    q183 --> q179
    q183 --> q181
    q184 --> q21
```

## RegexpTokenContent

```mermaid
flowchart TD
    q22(["RegexpTokenContent__Start (22)<br/>RuleStart"])
    q23(["RegexpTokenContent__Stop (23)<br/>RuleStop"])
    q185["RegexpTokenContent_Regexp_RegexLiteral (185)<br/>Basic<br/>"]
    q186["RegexpTokenContent__Basic (186)<br/>Basic<br/>"]

    q22 --> q185
    q185 -->|"tok(RegexLiteral)"| q186
    q186 --> q23
```

## KeywordTokenContent

```mermaid
flowchart TD
    q24(["KeywordTokenContent__Start (24)<br/>RuleStart"])
    q25(["KeywordTokenContent__Stop (25)<br/>RuleStop"])
    q187["KeywordTokenContent__Basic_0 (187)<br/>Basic<br/>"]
    q188["KeywordTokenContent__Basic_1 (188)<br/>Basic<br/>"]

    q24 --> q187
    q187 -.->|"[Keyword]"| q188
    q188 --> q25
```

## TokenCommand

```mermaid
flowchart TD
    q26(["TokenCommand__Start (26)<br/>RuleStart"])
    q27(["TokenCommand__Stop (27)<br/>RuleStop"])
    q189["TokenCommand_DashGreaterThan (189)<br/>Basic<br/>"]
    q190["TokenCommand_Type_push (190)<br/>Basic<br/>"]
    q191["TokenCommand__Basic_0 (191)<br/>Basic<br/>"]
    q192["TokenCommand_Type_pop (192)<br/>Basic<br/>"]
    q193["TokenCommand__Basic_1 (193)<br/>Basic<br/>"]
    q194["TokenCommand_Type_mode (194)<br/>Basic<br/>"]
    q195["TokenCommand__Basic_2 (195)<br/>Basic<br/>"]
    q196{"TokenCommand__Basic_3 (196)<br/>Basic<br/><br/>dec=15"}
    q197["TokenCommand__BlockEnd_0 (197)<br/>BlockEnd<br/>"]
    q198["TokenCommand_LeftParen (198)<br/>Basic<br/>"]
    q199["TokenCommand_Mode_ID (199)<br/>Basic<br/>"]
    q200["TokenCommand__Basic_4 (200)<br/>Basic<br/>"]
    q201["TokenCommand_Default_default (201)<br/>Basic<br/>"]
    q202["TokenCommand__Basic_5 (202)<br/>Basic<br/>"]
    q203{"TokenCommand__Basic_6 (203)<br/>Basic<br/><br/>dec=16"}
    q204["TokenCommand__BlockEnd_1 (204)<br/>BlockEnd<br/>"]
    q205["TokenCommand_RightParen (205)<br/>Basic<br/>"]
    q206["TokenCommand__Basic_7 (206)<br/>Basic<br/>"]
    q207{"TokenCommand__Basic_8 (207)<br/>Basic<br/><br/>dec=17"}

    q26 --> q189
    q189 -->|"tok(&quot;->&quot;)"| q196
    q190 -->|"tok(&quot;push&quot;)"| q191
    q191 --> q197
    q192 -->|"tok(&quot;pop&quot;)"| q193
    q193 --> q197
    q194 -->|"tok(&quot;mode&quot;)"| q195
    q195 --> q197
    q196 --> q190
    q196 --> q192
    q196 --> q194
    q197 --> q207
    q198 -->|"tok(&quot;(&quot;)"| q203
    q199 -->|"tok(ID)"| q200
    q200 --> q204
    q201 -->|"tok(&quot;default&quot;)"| q202
    q202 --> q204
    q203 --> q199
    q203 --> q201
    q204 --> q205
    q205 -->|"tok(&quot;)&quot;)"| q206
    q206 --> q27
    q207 --> q198
    q207 --> q206
```

## TokenGroup

```mermaid
flowchart TD
    q28(["TokenGroup__Start (28)<br/>RuleStart"])
    q29(["TokenGroup__Stop (29)<br/>RuleStop"])
    q208["TokenGroup_Modifier_TokenModifier (208)<br/>Basic<br/>"]
    q209["TokenGroup__Basic_0 (209)<br/>Basic<br/>"]
    q210{"TokenGroup__Basic_1 (210)<br/>Basic<br/><br/>dec=18"}
    q211["TokenGroup_token (211)<br/>Basic<br/>"]
    q212["TokenGroup_group (212)<br/>Basic<br/>"]
    q213["TokenGroup_Name_ID (213)<br/>Basic<br/>"]
    q214["TokenGroup_LeftBrace (214)<br/>Basic<br/>"]
    q215["TokenGroup_TokenRefs_ID (215)<br/>Basic<br/>"]
    q216["TokenGroup__Basic_2 (216)<br/>Basic<br/>"]
    q217["TokenGroup__Basic_3 (217)<br/>Basic<br/>"]
    q218["TokenGroup__Basic_4 (218)<br/>Basic<br/>"]
    q219["TokenGroup_keywords (219)<br/>Basic<br/>"]
    q220["TokenGroup_KeywordSelectors_RegexLiteral (220)<br/>Basic<br/>"]
    q221["TokenGroup__Basic_5 (221)<br/>Basic<br/>"]
    q222{"TokenGroup__Basic_6 (222)<br/>Basic<br/><br/>dec=19"}
    q223["TokenGroup__BlockEnd (223)<br/>BlockEnd<br/>"]
    q224{"TokenGroup__LoopEntry (224)<br/>LoopEntry<br/><br/>dec=20"}
    q225["TokenGroup__LoopEnd (225)<br/>LoopEnd<br/>"]
    q226["TokenGroup__LoopBack (226)<br/>LoopBack<br/>"]
    q227["TokenGroup_RightBrace (227)<br/>Basic<br/>"]
    q228["TokenGroup__Basic_7 (228)<br/>Basic<br/>"]
    q229["TokenGroup__Basic_8 (229)<br/>Basic<br/>"]
    q230{"TokenGroup__Basic_9 (230)<br/>Basic<br/><br/>dec=21"}
    q231["TokenGroup_Semicolon (231)<br/>Basic<br/>"]
    q232["TokenGroup__Basic_10 (232)<br/>Basic<br/>"]
    q233{"TokenGroup__Basic_11 (233)<br/>Basic<br/><br/>dec=22"}

    q28 --> q210
    q208 -->|"tok(TokenModifier)"| q209
    q209 --> q211
    q210 --> q208
    q210 --> q209
    q211 -->|"tok(&quot;token&quot;)"| q212
    q212 -->|"tok(&quot;group&quot;)"| q213
    q213 -->|"tok(ID)"| q214
    q214 -->|"tok(&quot;{&quot;)"| q224
    q215 -->|"tok(ID)"| q216
    q216 --> q223
    q217 -.->|"[Keyword]"| q218
    q218 --> q223
    q219 -->|"tok(&quot;keywords&quot;)"| q220
    q220 -->|"tok(RegexLiteral)"| q221
    q221 --> q223
    q222 --> q215
    q222 --> q217
    q222 --> q219
    q223 --> q226
    q224 --> q222
    q224 --> q225
    q225 --> q227
    q226 --> q224
    q227 -->|"tok(&quot;}&quot;)"| q230
    q228 -.->|"[TokenCommand]"| q229
    q229 --> q233
    q230 --> q228
    q230 --> q229
    q231 -->|"tok(&quot;;&quot;)"| q232
    q232 --> q29
    q233 --> q231
    q233 --> q232
```

## TokenMode

```mermaid
flowchart TD
    q30(["TokenMode__Start (30)<br/>RuleStart"])
    q31(["TokenMode__Stop (31)<br/>RuleStop"])
    q234["TokenMode_token (234)<br/>Basic<br/>"]
    q235["TokenMode_mode (235)<br/>Basic<br/>"]
    q236["TokenMode_Name_ID (236)<br/>Basic<br/>"]
    q237["TokenMode__Basic_0 (237)<br/>Basic<br/>"]
    q238["TokenMode_Default_default (238)<br/>Basic<br/>"]
    q239["TokenMode__Basic_1 (239)<br/>Basic<br/>"]
    q240{"TokenMode__Basic_2 (240)<br/>Basic<br/><br/>dec=23"}
    q241["TokenMode__BlockEnd (241)<br/>BlockEnd<br/>"]
    q242["TokenMode_LeftBrace (242)<br/>Basic<br/>"]
    q243["TokenMode__Basic_3 (243)<br/>Basic<br/>"]
    q244["TokenMode__Basic_4 (244)<br/>Basic<br/>"]
    q245{"TokenMode__LoopEntry (245)<br/>LoopEntry<br/><br/>dec=24"}
    q246["TokenMode__LoopEnd (246)<br/>LoopEnd<br/>"]
    q247["TokenMode__LoopBack (247)<br/>LoopBack<br/>"]
    q248["TokenMode_RightBrace (248)<br/>Basic<br/>"]
    q249["TokenMode__Basic_5 (249)<br/>Basic<br/>"]

    q30 --> q234
    q234 -->|"tok(&quot;token&quot;)"| q235
    q235 -->|"tok(&quot;mode&quot;)"| q240
    q236 -->|"tok(ID)"| q237
    q237 --> q241
    q238 -->|"tok(&quot;default&quot;)"| q239
    q239 --> q241
    q240 --> q236
    q240 --> q238
    q241 --> q242
    q242 -->|"tok(&quot;{&quot;)"| q245
    q243 -.->|"[TokenModeMember]"| q244
    q244 --> q247
    q245 --> q243
    q245 --> q246
    q246 --> q248
    q247 --> q245
    q248 -->|"tok(&quot;}&quot;)"| q249
    q249 --> q31
```

## TokenModeMember

```mermaid
flowchart TD
    q32(["TokenModeMember__Start (32)<br/>RuleStart"])
    q33(["TokenModeMember__Stop (33)<br/>RuleStop"])
    q250["TokenModeMember__Basic_0 (250)<br/>Basic<br/>"]
    q251["TokenModeMember__Basic_1 (251)<br/>Basic<br/>"]
    q252["TokenModeMember__Basic_2 (252)<br/>Basic<br/>"]
    q253["TokenModeMember__Basic_3 (253)<br/>Basic<br/>"]
    q254["TokenModeMember__Basic_4 (254)<br/>Basic<br/>"]
    q255["TokenModeMember__Basic_5 (255)<br/>Basic<br/>"]
    q256["TokenModeMember__Basic_6 (256)<br/>Basic<br/>"]
    q257["TokenModeMember__Basic_7 (257)<br/>Basic<br/>"]
    q258["TokenModeMember__Basic_8 (258)<br/>Basic<br/>"]
    q259["TokenModeMember__Basic_9 (259)<br/>Basic<br/>"]
    q260{"TokenModeMember__Basic_10 (260)<br/>Basic<br/><br/>dec=25"}
    q261["TokenModeMember__BlockEnd (261)<br/>BlockEnd<br/>"]

    q32 --> q260
    q250 -.->|"[TokenDeclUsage]"| q251
    q251 --> q261
    q252 -.->|"[TokenGroupUsage]"| q253
    q253 --> q261
    q254 -.->|"[TokenUsage]"| q255
    q255 --> q261
    q256 -.->|"[KeywordUsage]"| q257
    q257 --> q261
    q258 -.->|"[KeywordSelector]"| q259
    q259 --> q261
    q260 --> q250
    q260 --> q252
    q260 --> q254
    q260 --> q256
    q260 --> q258
    q261 --> q33
```

## TokenDeclUsage

```mermaid
flowchart TD
    q34(["TokenDeclUsage__Start (34)<br/>RuleStart"])
    q35(["TokenDeclUsage__Stop (35)<br/>RuleStop"])
    q262["TokenDeclUsage__Basic_0 (262)<br/>Basic<br/>"]
    q263["TokenDeclUsage__Basic_1 (263)<br/>Basic<br/>"]

    q34 --> q262
    q262 -.->|"[TokenDecl]"| q263
    q263 --> q35
```

## TokenGroupUsage

```mermaid
flowchart TD
    q36(["TokenGroupUsage__Start (36)<br/>RuleStart"])
    q37(["TokenGroupUsage__Stop (37)<br/>RuleStop"])
    q264["TokenGroupUsage__Basic_0 (264)<br/>Basic<br/>"]
    q265["TokenGroupUsage__Basic_1 (265)<br/>Basic<br/>"]

    q36 --> q264
    q264 -.->|"[TokenGroup]"| q265
    q265 --> q37
```

## TokenUsage

```mermaid
flowchart TD
    q38(["TokenUsage__Start (38)<br/>RuleStart"])
    q39(["TokenUsage__Stop (39)<br/>RuleStop"])
    q266["TokenUsage_Modifier_TokenModifier (266)<br/>Basic<br/>"]
    q267["TokenUsage__Basic_0 (267)<br/>Basic<br/>"]
    q268{"TokenUsage__Basic_1 (268)<br/>Basic<br/><br/>dec=26"}
    q269["TokenUsage_TokenRef_ID (269)<br/>Basic<br/>"]
    q270["TokenUsage__Basic_2 (270)<br/>Basic<br/>"]
    q271["TokenUsage__Basic_3 (271)<br/>Basic<br/>"]
    q272{"TokenUsage__Basic_4 (272)<br/>Basic<br/><br/>dec=27"}
    q273["TokenUsage_Semicolon (273)<br/>Basic<br/>"]
    q274["TokenUsage__Basic_5 (274)<br/>Basic<br/>"]
    q275{"TokenUsage__Basic_6 (275)<br/>Basic<br/><br/>dec=28"}

    q38 --> q268
    q266 -->|"tok(TokenModifier)"| q267
    q267 --> q269
    q268 --> q266
    q268 --> q267
    q269 -->|"tok(ID)"| q272
    q270 -.->|"[TokenCommand]"| q271
    q271 --> q275
    q272 --> q270
    q272 --> q271
    q273 -->|"tok(&quot;;&quot;)"| q274
    q274 --> q39
    q275 --> q273
    q275 --> q274
```

## KeywordUsage

```mermaid
flowchart TD
    q40(["KeywordUsage__Start (40)<br/>RuleStart"])
    q41(["KeywordUsage__Stop (41)<br/>RuleStop"])
    q276["KeywordUsage_Modifier_TokenModifier (276)<br/>Basic<br/>"]
    q277["KeywordUsage__Basic_0 (277)<br/>Basic<br/>"]
    q278{"KeywordUsage__Basic_1 (278)<br/>Basic<br/><br/>dec=29"}
    q279["KeywordUsage__Basic_2 (279)<br/>Basic<br/>"]
    q280["KeywordUsage__Basic_3 (280)<br/>Basic<br/>"]
    q281["KeywordUsage__Basic_4 (281)<br/>Basic<br/>"]
    q282{"KeywordUsage__Basic_5 (282)<br/>Basic<br/><br/>dec=30"}
    q283["KeywordUsage_Semicolon (283)<br/>Basic<br/>"]
    q284["KeywordUsage__Basic_6 (284)<br/>Basic<br/>"]
    q285{"KeywordUsage__Basic_7 (285)<br/>Basic<br/><br/>dec=31"}

    q40 --> q278
    q276 -->|"tok(TokenModifier)"| q277
    q277 --> q279
    q278 --> q276
    q278 --> q277
    q279 -.->|"[Keyword]"| q282
    q280 -.->|"[TokenCommand]"| q281
    q281 --> q285
    q282 --> q280
    q282 --> q281
    q283 -->|"tok(&quot;;&quot;)"| q284
    q284 --> q41
    q285 --> q283
    q285 --> q284
```

## KeywordSelector

```mermaid
flowchart TD
    q42(["KeywordSelector__Start (42)<br/>RuleStart"])
    q43(["KeywordSelector__Stop (43)<br/>RuleStop"])
    q286["KeywordSelector_keywords (286)<br/>Basic<br/>"]
    q287["KeywordSelector_Selector_RegexLiteral (287)<br/>Basic<br/>"]
    q288["KeywordSelector_Semicolon (288)<br/>Basic<br/>"]
    q289["KeywordSelector__Basic_0 (289)<br/>Basic<br/>"]
    q290{"KeywordSelector__Basic_1 (290)<br/>Basic<br/><br/>dec=32"}

    q42 --> q286
    q286 -->|"tok(&quot;keywords&quot;)"| q287
    q287 -->|"tok(RegexLiteral)"| q290
    q288 -->|"tok(&quot;;&quot;)"| q289
    q289 --> q43
    q290 --> q288
    q290 --> q289
```

## Alternatives

```mermaid
flowchart TD
    q44(["Alternatives__Start (44)<br/>RuleStart"])
    q45(["Alternatives__Stop (45)<br/>RuleStop"])
    q291["Alternatives__Basic_0 (291)<br/>Basic<br/>"]
    q292["Alternatives_Pipe (292)<br/>Basic<br/>"]
    q293["Alternatives__Basic_1 (293)<br/>Basic<br/>"]
    q294["Alternatives__Basic_2 (294)<br/>Basic<br/>"]
    q295{"Alternatives__LoopBack (295)<br/>LoopBack<br/><br/>dec=33"}
    q296["Alternatives__LoopEnd (296)<br/>LoopEnd<br/>"]
    q297{"Alternatives__Basic_3 (297)<br/>Basic<br/><br/>dec=34"}

    q44 --> q291
    q291 -.->|"[Group]"| q297
    q292 -->|"tok(&quot;|&quot;)"| q293
    q293 -.->|"[Group]"| q294
    q294 --> q295
    q295 --> q292
    q295 --> q296
    q296 --> q45
    q297 --> q292
    q297 --> q296
```

## Group

```mermaid
flowchart TD
    q46(["Group__Start (46)<br/>RuleStart"])
    q47(["Group__Stop (47)<br/>RuleStop"])
    q298["Group__Basic_0 (298)<br/>Basic<br/>"]
    q299["Group__Basic_1 (299)<br/>Basic<br/>"]
    q300["Group__Basic_2 (300)<br/>Basic<br/>"]
    q301{"Group__LoopBack (301)<br/>LoopBack<br/><br/>dec=35"}
    q302["Group__LoopEnd (302)<br/>LoopEnd<br/>"]
    q303{"Group__Basic_3 (303)<br/>Basic<br/><br/>dec=36"}

    q46 --> q298
    q298 -.->|"[Element]"| q303
    q299 -.->|"[Element]"| q300
    q300 --> q301
    q301 --> q299
    q301 --> q302
    q302 --> q47
    q303 --> q299
    q303 --> q302
```

## Element

```mermaid
flowchart TD
    q48(["Element__Start (48)<br/>RuleStart"])
    q49(["Element__Stop (49)<br/>RuleStop"])
    q304["Element__Basic_0 (304)<br/>Basic<br/>"]
    q305["Element__Basic_1 (305)<br/>Basic<br/>"]
    q306["Element__Basic_2 (306)<br/>Basic<br/>"]
    q307["Element__Basic_3 (307)<br/>Basic<br/>"]
    q308["Element__Basic_4 (308)<br/>Basic<br/>"]
    q309["Element__Basic_5 (309)<br/>Basic<br/>"]
    q310["Element__Basic_6 (310)<br/>Basic<br/>"]
    q311["Element__Basic_7 (311)<br/>Basic<br/>"]
    q312["Element_LeftParen (312)<br/>Basic<br/>"]
    q313["Element__Basic_8 (313)<br/>Basic<br/>"]
    q314["Element_RightParen (314)<br/>Basic<br/>"]
    q315["Element__Basic_9 (315)<br/>Basic<br/>"]
    q316{"Element__Basic_10 (316)<br/>Basic<br/><br/>dec=37"}
    q317["Element__BlockEnd (317)<br/>BlockEnd<br/>"]
    q318["Element_Cardinality_Cardinality (318)<br/>Basic<br/>"]
    q319["Element__Basic_11 (319)<br/>Basic<br/>"]
    q320{"Element__Basic_12 (320)<br/>Basic<br/><br/>dec=38"}

    q48 --> q316
    q304 -.->|"[Keyword]"| q305
    q305 --> q317
    q306 -.->|"[Assignment]"| q307
    q307 --> q317
    q308 -.->|"[RuleCall]"| q309
    q309 --> q317
    q310 -.->|"[Action]"| q311
    q311 --> q317
    q312 -->|"tok(&quot;(&quot;)"| q313
    q313 -.->|"[Alternatives]"| q314
    q314 -->|"tok(&quot;)&quot;)"| q315
    q315 --> q317
    q316 --> q304
    q316 --> q306
    q316 --> q308
    q316 --> q310
    q316 --> q312
    q317 --> q320
    q318 -->|"tok(Cardinality)"| q319
    q319 --> q49
    q320 --> q318
    q320 --> q319
```

## Keyword

```mermaid
flowchart TD
    q50(["Keyword__Start (50)<br/>RuleStart"])
    q51(["Keyword__Stop (51)<br/>RuleStop"])
    q321["Keyword_Value_StringLiteral (321)<br/>Basic<br/>"]
    q322["Keyword__Basic (322)<br/>Basic<br/>"]

    q50 --> q321
    q321 -->|"tok(StringLiteral)"| q322
    q322 --> q51
```

## Assignment

```mermaid
flowchart TD
    q52(["Assignment__Start (52)<br/>RuleStart"])
    q53(["Assignment__Stop (53)<br/>RuleStop"])
    q323["Assignment_Property_ID (323)<br/>Basic<br/>"]
    q324["Assignment_Operator_PlusEquals (324)<br/>Basic<br/>"]
    q325["Assignment__Basic_0 (325)<br/>Basic<br/>"]
    q326["Assignment_Operator_Equals (326)<br/>Basic<br/>"]
    q327["Assignment__Basic_1 (327)<br/>Basic<br/>"]
    q328["Assignment_Operator_QuestionEquals (328)<br/>Basic<br/>"]
    q329["Assignment__Basic_2 (329)<br/>Basic<br/>"]
    q330{"Assignment__Basic_3 (330)<br/>Basic<br/><br/>dec=39"}
    q331["Assignment__BlockEnd (331)<br/>BlockEnd<br/>"]
    q332["Assignment__Basic_4 (332)<br/>Basic<br/>"]
    q333["Assignment__Basic_5 (333)<br/>Basic<br/>"]

    q52 --> q323
    q323 -->|"tok(ID)"| q330
    q324 -->|"tok(&quot;+=&quot;)"| q325
    q325 --> q331
    q326 -->|"tok(&quot;=&quot;)"| q327
    q327 --> q331
    q328 -->|"tok(&quot;?=&quot;)"| q329
    q329 --> q331
    q330 --> q324
    q330 --> q326
    q330 --> q328
    q331 --> q332
    q332 -.->|"[Assignable]"| q333
    q333 --> q53
```

## Assignable

```mermaid
flowchart TD
    q54(["Assignable__Start (54)<br/>RuleStart"])
    q55(["Assignable__Stop (55)<br/>RuleStop"])
    q334["Assignable__Basic_0 (334)<br/>Basic<br/>"]
    q335["Assignable__Basic_1 (335)<br/>Basic<br/>"]
    q336["Assignable__Basic_2 (336)<br/>Basic<br/>"]
    q337["Assignable__Basic_3 (337)<br/>Basic<br/>"]
    q338["Assignable__Basic_4 (338)<br/>Basic<br/>"]
    q339["Assignable__Basic_5 (339)<br/>Basic<br/>"]
    q340["Assignable_LeftParen (340)<br/>Basic<br/>"]
    q341["Assignable__Basic_6 (341)<br/>Basic<br/>"]
    q342["Assignable_RightParen (342)<br/>Basic<br/>"]
    q343["Assignable__Basic_7 (343)<br/>Basic<br/>"]
    q344{"Assignable__Basic_8 (344)<br/>Basic<br/><br/>dec=40"}
    q345["Assignable__BlockEnd (345)<br/>BlockEnd<br/>"]

    q54 --> q344
    q334 -.->|"[Keyword]"| q335
    q335 --> q345
    q336 -.->|"[RuleCall]"| q337
    q337 --> q345
    q338 -.->|"[CrossRef]"| q339
    q339 --> q345
    q340 -->|"tok(&quot;(&quot;)"| q341
    q341 -.->|"[AssignableAlternatives]"| q342
    q342 -->|"tok(&quot;)&quot;)"| q343
    q343 --> q345
    q344 --> q334
    q344 --> q336
    q344 --> q338
    q344 --> q340
    q345 --> q55
```

## AssignableWithoutAlts

```mermaid
flowchart TD
    q56(["AssignableWithoutAlts__Start (56)<br/>RuleStart"])
    q57(["AssignableWithoutAlts__Stop (57)<br/>RuleStop"])
    q346["AssignableWithoutAlts__Basic_0 (346)<br/>Basic<br/>"]
    q347["AssignableWithoutAlts__Basic_1 (347)<br/>Basic<br/>"]
    q348["AssignableWithoutAlts__Basic_2 (348)<br/>Basic<br/>"]
    q349["AssignableWithoutAlts__Basic_3 (349)<br/>Basic<br/>"]
    q350["AssignableWithoutAlts__Basic_4 (350)<br/>Basic<br/>"]
    q351["AssignableWithoutAlts__Basic_5 (351)<br/>Basic<br/>"]
    q352{"AssignableWithoutAlts__Basic_6 (352)<br/>Basic<br/><br/>dec=41"}
    q353["AssignableWithoutAlts__BlockEnd (353)<br/>BlockEnd<br/>"]

    q56 --> q352
    q346 -.->|"[Keyword]"| q347
    q347 --> q353
    q348 -.->|"[RuleCall]"| q349
    q349 --> q353
    q350 -.->|"[CrossRef]"| q351
    q351 --> q353
    q352 --> q346
    q352 --> q348
    q352 --> q350
    q353 --> q57
```

## AssignableAlternatives

```mermaid
flowchart TD
    q58(["AssignableAlternatives__Start (58)<br/>RuleStart"])
    q59(["AssignableAlternatives__Stop (59)<br/>RuleStop"])
    q354["AssignableAlternatives__Basic_0 (354)<br/>Basic<br/>"]
    q355["AssignableAlternatives_Pipe (355)<br/>Basic<br/>"]
    q356["AssignableAlternatives__Basic_1 (356)<br/>Basic<br/>"]
    q357["AssignableAlternatives__Basic_2 (357)<br/>Basic<br/>"]
    q358{"AssignableAlternatives__LoopBack (358)<br/>LoopBack<br/><br/>dec=42"}
    q359["AssignableAlternatives__LoopEnd (359)<br/>LoopEnd<br/>"]
    q360{"AssignableAlternatives__Basic_3 (360)<br/>Basic<br/><br/>dec=43"}

    q58 --> q354
    q354 -.->|"[AssignableWithoutAlts]"| q360
    q355 -->|"tok(&quot;|&quot;)"| q356
    q356 -.->|"[AssignableWithoutAlts]"| q357
    q357 --> q358
    q358 --> q355
    q358 --> q359
    q359 --> q59
    q360 --> q355
    q360 --> q359
```

## CrossRef

```mermaid
flowchart TD
    q60(["CrossRef__Start (60)<br/>RuleStart"])
    q61(["CrossRef__Stop (61)<br/>RuleStop"])
    q361["CrossRef_LeftBracket (361)<br/>Basic<br/>"]
    q362["CrossRef_Type_ID (362)<br/>Basic<br/>"]
    q363["CrossRef_Colon (363)<br/>Basic<br/>"]
    q364["CrossRef__Basic_0 (364)<br/>Basic<br/>"]
    q365["CrossRef__Basic_1 (365)<br/>Basic<br/>"]
    q366{"CrossRef__Basic_2 (366)<br/>Basic<br/><br/>dec=44"}
    q367["CrossRef_RightBracket (367)<br/>Basic<br/>"]
    q368["CrossRef__Basic_3 (368)<br/>Basic<br/>"]

    q60 --> q361
    q361 -->|"tok(&quot;[&quot;)"| q362
    q362 -->|"tok(ID)"| q366
    q363 -->|"tok(&quot;:&quot;)"| q364
    q364 -.->|"[RuleCall]"| q365
    q365 --> q367
    q366 --> q363
    q366 --> q365
    q367 -->|"tok(&quot;]&quot;)"| q368
    q368 --> q61
```

## RuleCall

```mermaid
flowchart TD
    q62(["RuleCall__Start (62)<br/>RuleStart"])
    q63(["RuleCall__Stop (63)<br/>RuleStop"])
    q369["RuleCall_Rule_ID (369)<br/>Basic<br/>"]
    q370["RuleCall__Basic (370)<br/>Basic<br/>"]

    q62 --> q369
    q369 -->|"tok(ID)"| q370
    q370 --> q63
```

## Action

```mermaid
flowchart TD
    q64(["Action__Start (64)<br/>RuleStart"])
    q65(["Action__Stop (65)<br/>RuleStop"])
    q371["Action_LeftBrace (371)<br/>Basic<br/>"]
    q372["Action_Type_ID (372)<br/>Basic<br/>"]
    q373["Action_Dot (373)<br/>Basic<br/>"]
    q374["Action_Property_ID (374)<br/>Basic<br/>"]
    q375["Action_Operator_PlusEquals (375)<br/>Basic<br/>"]
    q376["Action__Basic_0 (376)<br/>Basic<br/>"]
    q377["Action_Operator_Equals (377)<br/>Basic<br/>"]
    q378["Action__Basic_1 (378)<br/>Basic<br/>"]
    q379{"Action__Basic_2 (379)<br/>Basic<br/><br/>dec=45"}
    q380["Action__BlockEnd (380)<br/>BlockEnd<br/>"]
    q381["Action_current (381)<br/>Basic<br/>"]
    q382["Action__Basic_3 (382)<br/>Basic<br/>"]
    q383{"Action__Basic_4 (383)<br/>Basic<br/><br/>dec=46"}
    q384["Action_RightBrace (384)<br/>Basic<br/>"]
    q385["Action__Basic_5 (385)<br/>Basic<br/>"]

    q64 --> q371
    q371 -->|"tok(&quot;{&quot;)"| q372
    q372 -->|"tok(ID)"| q383
    q373 -->|"tok(&quot;.&quot;)"| q374
    q374 -->|"tok(ID)"| q379
    q375 -->|"tok(&quot;+=&quot;)"| q376
    q376 --> q380
    q377 -->|"tok(&quot;=&quot;)"| q378
    q378 --> q380
    q379 --> q375
    q379 --> q377
    q380 --> q381
    q381 -->|"tok(&quot;current&quot;)"| q382
    q382 --> q384
    q383 --> q373
    q383 --> q382
    q384 -->|"tok(&quot;}&quot;)"| q385
    q385 --> q65
```

## CompositeRule

```mermaid
flowchart TD
    q66(["CompositeRule__Start (66)<br/>RuleStart"])
    q67(["CompositeRule__Stop (67)<br/>RuleStop"])
    q386["CompositeRule_composite (386)<br/>Basic<br/>"]
    q387["CompositeRule_Name_ID (387)<br/>Basic<br/>"]
    q388["CompositeRule_Colon (388)<br/>Basic<br/>"]
    q389["CompositeRule__Basic_0 (389)<br/>Basic<br/>"]
    q390["CompositeRule_Semicolon (390)<br/>Basic<br/>"]
    q391["CompositeRule__Basic_1 (391)<br/>Basic<br/>"]
    q392{"CompositeRule__Basic_2 (392)<br/>Basic<br/><br/>dec=47"}

    q66 --> q386
    q386 -->|"tok(&quot;composite&quot;)"| q387
    q387 -->|"tok(ID)"| q388
    q388 -->|"tok(&quot;:&quot;)"| q389
    q389 -.->|"[CompositeAlternatives]"| q392
    q390 -->|"tok(&quot;;&quot;)"| q391
    q391 --> q67
    q392 --> q390
    q392 --> q391
```

## CompositeAlternatives

```mermaid
flowchart TD
    q68(["CompositeAlternatives__Start (68)<br/>RuleStart"])
    q69(["CompositeAlternatives__Stop (69)<br/>RuleStop"])
    q393["CompositeAlternatives__Basic_0 (393)<br/>Basic<br/>"]
    q394["CompositeAlternatives_Pipe (394)<br/>Basic<br/>"]
    q395["CompositeAlternatives__Basic_1 (395)<br/>Basic<br/>"]
    q396["CompositeAlternatives__Basic_2 (396)<br/>Basic<br/>"]
    q397{"CompositeAlternatives__LoopBack (397)<br/>LoopBack<br/><br/>dec=48"}
    q398["CompositeAlternatives__LoopEnd (398)<br/>LoopEnd<br/>"]
    q399{"CompositeAlternatives__Basic_3 (399)<br/>Basic<br/><br/>dec=49"}

    q68 --> q393
    q393 -.->|"[CompositeGroup]"| q399
    q394 -->|"tok(&quot;|&quot;)"| q395
    q395 -.->|"[CompositeGroup]"| q396
    q396 --> q397
    q397 --> q394
    q397 --> q398
    q398 --> q69
    q399 --> q394
    q399 --> q398
```

## CompositeGroup

```mermaid
flowchart TD
    q70(["CompositeGroup__Start (70)<br/>RuleStart"])
    q71(["CompositeGroup__Stop (71)<br/>RuleStop"])
    q400["CompositeGroup__Basic_0 (400)<br/>Basic<br/>"]
    q401["CompositeGroup__Basic_1 (401)<br/>Basic<br/>"]
    q402["CompositeGroup__Basic_2 (402)<br/>Basic<br/>"]
    q403{"CompositeGroup__LoopBack (403)<br/>LoopBack<br/><br/>dec=50"}
    q404["CompositeGroup__LoopEnd (404)<br/>LoopEnd<br/>"]
    q405{"CompositeGroup__Basic_3 (405)<br/>Basic<br/><br/>dec=51"}

    q70 --> q400
    q400 -.->|"[CompositeElement]"| q405
    q401 -.->|"[CompositeElement]"| q402
    q402 --> q403
    q403 --> q401
    q403 --> q404
    q404 --> q71
    q405 --> q401
    q405 --> q404
```

## CompositeElement

```mermaid
flowchart TD
    q72(["CompositeElement__Start (72)<br/>RuleStart"])
    q73(["CompositeElement__Stop (73)<br/>RuleStop"])
    q406["CompositeElement__Basic_0 (406)<br/>Basic<br/>"]
    q407["CompositeElement__Basic_1 (407)<br/>Basic<br/>"]
    q408["CompositeElement__Basic_2 (408)<br/>Basic<br/>"]
    q409["CompositeElement__Basic_3 (409)<br/>Basic<br/>"]
    q410["CompositeElement_LeftParen (410)<br/>Basic<br/>"]
    q411["CompositeElement__Basic_4 (411)<br/>Basic<br/>"]
    q412["CompositeElement_RightParen (412)<br/>Basic<br/>"]
    q413["CompositeElement__Basic_5 (413)<br/>Basic<br/>"]
    q414{"CompositeElement__Basic_6 (414)<br/>Basic<br/><br/>dec=52"}
    q415["CompositeElement__BlockEnd (415)<br/>BlockEnd<br/>"]
    q416["CompositeElement_Cardinality_Cardinality (416)<br/>Basic<br/>"]
    q417["CompositeElement__Basic_7 (417)<br/>Basic<br/>"]
    q418{"CompositeElement__Basic_8 (418)<br/>Basic<br/><br/>dec=53"}

    q72 --> q414
    q406 -.->|"[Keyword]"| q407
    q407 --> q415
    q408 -.->|"[RuleCall]"| q409
    q409 --> q415
    q410 -->|"tok(&quot;(&quot;)"| q411
    q411 -.->|"[CompositeAlternatives]"| q412
    q412 -->|"tok(&quot;)&quot;)"| q413
    q413 --> q415
    q414 --> q406
    q414 --> q408
    q414 --> q410
    q415 --> q418
    q416 -->|"tok(Cardinality)"| q417
    q417 --> q73
    q418 --> q416
    q418 --> q417
```

## InfixRule

```mermaid
flowchart TD
    q74(["InfixRule__Start (74)<br/>RuleStart"])
    q75(["InfixRule__Stop (75)<br/>RuleStop"])
    q419["InfixRule_infix (419)<br/>Basic<br/>"]
    q420["InfixRule_Name_ID (420)<br/>Basic<br/>"]
    q421["InfixRule_on (421)<br/>Basic<br/>"]
    q422["InfixRule__Basic_0 (422)<br/>Basic<br/>"]
    q423["InfixRule_returns (423)<br/>Basic<br/>"]
    q424["InfixRule_ReturnType_ID (424)<br/>Basic<br/>"]
    q425["InfixRule__Basic_1 (425)<br/>Basic<br/>"]
    q426{"InfixRule__Basic_2 (426)<br/>Basic<br/><br/>dec=54"}
    q427["InfixRule_Colon (427)<br/>Basic<br/>"]
    q428["InfixRule__Basic_3 (428)<br/>Basic<br/>"]
    q429["InfixRule_GreaterThan (429)<br/>Basic<br/>"]
    q430["InfixRule__Basic_4 (430)<br/>Basic<br/>"]
    q431["InfixRule__Basic_5 (431)<br/>Basic<br/>"]
    q432{"InfixRule__LoopEntry (432)<br/>LoopEntry<br/><br/>dec=55"}
    q433["InfixRule__LoopEnd (433)<br/>LoopEnd<br/>"]
    q434["InfixRule__LoopBack (434)<br/>LoopBack<br/>"]
    q435["InfixRule_Semicolon (435)<br/>Basic<br/>"]
    q436["InfixRule__Basic_6 (436)<br/>Basic<br/>"]
    q437{"InfixRule__Basic_7 (437)<br/>Basic<br/><br/>dec=56"}

    q74 --> q419
    q419 -->|"tok(&quot;infix&quot;)"| q420
    q420 -->|"tok(ID)"| q421
    q421 -->|"tok(&quot;on&quot;)"| q422
    q422 -.->|"[RuleCall]"| q426
    q423 -->|"tok(&quot;returns&quot;)"| q424
    q424 -->|"tok(ID)"| q425
    q425 --> q427
    q426 --> q423
    q426 --> q425
    q427 -->|"tok(&quot;:&quot;)"| q428
    q428 -.->|"[PrecedenceGroup]"| q432
    q429 -->|"tok(&quot;>&quot;)"| q430
    q430 -.->|"[PrecedenceGroup]"| q431
    q431 --> q434
    q432 --> q429
    q432 --> q433
    q433 --> q437
    q434 --> q432
    q435 -->|"tok(&quot;;&quot;)"| q436
    q436 --> q75
    q437 --> q435
    q437 --> q436
```

## PrecedenceGroup

```mermaid
flowchart TD
    q76(["PrecedenceGroup__Start (76)<br/>RuleStart"])
    q77(["PrecedenceGroup__Stop (77)<br/>RuleStop"])
    q438["PrecedenceGroup_Associativity_left (438)<br/>Basic<br/>"]
    q439["PrecedenceGroup__Basic_0 (439)<br/>Basic<br/>"]
    q440["PrecedenceGroup_Associativity_right (440)<br/>Basic<br/>"]
    q441["PrecedenceGroup__Basic_1 (441)<br/>Basic<br/>"]
    q442{"PrecedenceGroup__Basic_2 (442)<br/>Basic<br/><br/>dec=57"}
    q443["PrecedenceGroup__BlockEnd (443)<br/>BlockEnd<br/>"]
    q444{"PrecedenceGroup__Basic_3 (444)<br/>Basic<br/><br/>dec=58"}
    q445["PrecedenceGroup__Basic_4 (445)<br/>Basic<br/>"]
    q446["PrecedenceGroup_Pipe (446)<br/>Basic<br/>"]
    q447["PrecedenceGroup__Basic_5 (447)<br/>Basic<br/>"]
    q448["PrecedenceGroup__Basic_6 (448)<br/>Basic<br/>"]
    q449{"PrecedenceGroup__LoopEntry (449)<br/>LoopEntry<br/><br/>dec=59"}
    q450["PrecedenceGroup__LoopEnd (450)<br/>LoopEnd<br/>"]
    q451["PrecedenceGroup__LoopBack (451)<br/>LoopBack<br/>"]

    q76 --> q444
    q438 -->|"tok(&quot;left&quot;)"| q439
    q439 --> q443
    q440 -->|"tok(&quot;right&quot;)"| q441
    q441 --> q443
    q442 --> q438
    q442 --> q440
    q443 --> q445
    q444 --> q442
    q444 --> q443
    q445 -.->|"[InfixOperator]"| q449
    q446 -->|"tok(&quot;|&quot;)"| q447
    q447 -.->|"[InfixOperator]"| q448
    q448 --> q451
    q449 --> q446
    q449 --> q450
    q450 --> q77
    q451 --> q449
```

## InfixOperator

```mermaid
flowchart TD
    q78(["InfixOperator__Start (78)<br/>RuleStart"])
    q79(["InfixOperator__Stop (79)<br/>RuleStop"])
    q452["InfixOperator__Basic_0 (452)<br/>Basic<br/>"]
    q453["InfixOperator__Basic_1 (453)<br/>Basic<br/>"]
    q454["InfixOperator__Basic_2 (454)<br/>Basic<br/>"]
    q455["InfixOperator__Basic_3 (455)<br/>Basic<br/>"]
    q456{"InfixOperator__Basic_4 (456)<br/>Basic<br/><br/>dec=60"}
    q457["InfixOperator__BlockEnd (457)<br/>BlockEnd<br/>"]

    q78 --> q456
    q452 -.->|"[Keyword]"| q453
    q453 --> q457
    q454 -.->|"[RuleCall]"| q455
    q455 --> q457
    q456 --> q452
    q456 --> q454
    q457 --> q79
```

