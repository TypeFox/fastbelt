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
    q197["TokenCommand__BlockEnd (197)<br/>BlockEnd<br/>"]
    q198["TokenCommand_LeftParen (198)<br/>Basic<br/>"]
    q199["TokenCommand_Mode_ID (199)<br/>Basic<br/>"]
    q200["TokenCommand_RightParen (200)<br/>Basic<br/>"]
    q201["TokenCommand__Basic_4 (201)<br/>Basic<br/>"]
    q202{"TokenCommand__Basic_5 (202)<br/>Basic<br/><br/>dec=16"}

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
    q197 --> q202
    q198 -->|"tok(&quot;(&quot;)"| q199
    q199 -->|"tok(ID)"| q200
    q200 -->|"tok(&quot;)&quot;)"| q201
    q201 --> q27
    q202 --> q198
    q202 --> q201
```

## TokenGroup

```mermaid
flowchart TD
    q28(["TokenGroup__Start (28)<br/>RuleStart"])
    q29(["TokenGroup__Stop (29)<br/>RuleStop"])
    q203["TokenGroup_token (203)<br/>Basic<br/>"]
    q204["TokenGroup_group (204)<br/>Basic<br/>"]
    q205["TokenGroup_Name_ID (205)<br/>Basic<br/>"]
    q206["TokenGroup_LeftBrace (206)<br/>Basic<br/>"]
    q207["TokenGroup_TokenRefs_ID (207)<br/>Basic<br/>"]
    q208["TokenGroup__Basic_0 (208)<br/>Basic<br/>"]
    q209["TokenGroup__Basic_1 (209)<br/>Basic<br/>"]
    q210["TokenGroup__Basic_2 (210)<br/>Basic<br/>"]
    q211["TokenGroup_keywords (211)<br/>Basic<br/>"]
    q212["TokenGroup_KeywordSelectors_RegexLiteral (212)<br/>Basic<br/>"]
    q213["TokenGroup__Basic_3 (213)<br/>Basic<br/>"]
    q214{"TokenGroup__Basic_4 (214)<br/>Basic<br/><br/>dec=17"}
    q215["TokenGroup__BlockEnd (215)<br/>BlockEnd<br/>"]
    q216{"TokenGroup__LoopEntry (216)<br/>LoopEntry<br/><br/>dec=18"}
    q217["TokenGroup__LoopEnd (217)<br/>LoopEnd<br/>"]
    q218["TokenGroup__LoopBack (218)<br/>LoopBack<br/>"]
    q219["TokenGroup_RightBrace (219)<br/>Basic<br/>"]
    q220["TokenGroup_Semicolon (220)<br/>Basic<br/>"]
    q221["TokenGroup__Basic_5 (221)<br/>Basic<br/>"]
    q222{"TokenGroup__Basic_6 (222)<br/>Basic<br/><br/>dec=19"}

    q28 --> q203
    q203 -->|"tok(&quot;token&quot;)"| q204
    q204 -->|"tok(&quot;group&quot;)"| q205
    q205 -->|"tok(ID)"| q206
    q206 -->|"tok(&quot;{&quot;)"| q216
    q207 -->|"tok(ID)"| q208
    q208 --> q215
    q209 -.->|"[Keyword]"| q210
    q210 --> q215
    q211 -->|"tok(&quot;keywords&quot;)"| q212
    q212 -->|"tok(RegexLiteral)"| q213
    q213 --> q215
    q214 --> q207
    q214 --> q209
    q214 --> q211
    q215 --> q218
    q216 --> q214
    q216 --> q217
    q217 --> q219
    q218 --> q216
    q219 -->|"tok(&quot;}&quot;)"| q222
    q220 -->|"tok(&quot;;&quot;)"| q221
    q221 --> q29
    q222 --> q220
    q222 --> q221
```

## TokenMode

```mermaid
flowchart TD
    q30(["TokenMode__Start (30)<br/>RuleStart"])
    q31(["TokenMode__Stop (31)<br/>RuleStop"])
    q223["TokenMode_token (223)<br/>Basic<br/>"]
    q224["TokenMode_mode (224)<br/>Basic<br/>"]
    q225["TokenMode_Name_ID (225)<br/>Basic<br/>"]
    q226["TokenMode_LeftBrace (226)<br/>Basic<br/>"]
    q227["TokenMode__Basic_0 (227)<br/>Basic<br/>"]
    q228["TokenMode__Basic_1 (228)<br/>Basic<br/>"]
    q229{"TokenMode__LoopEntry (229)<br/>LoopEntry<br/><br/>dec=20"}
    q230["TokenMode__LoopEnd (230)<br/>LoopEnd<br/>"]
    q231["TokenMode__LoopBack (231)<br/>LoopBack<br/>"]
    q232["TokenMode_RightBrace (232)<br/>Basic<br/>"]
    q233["TokenMode__Basic_2 (233)<br/>Basic<br/>"]

    q30 --> q223
    q223 -->|"tok(&quot;token&quot;)"| q224
    q224 -->|"tok(&quot;mode&quot;)"| q225
    q225 -->|"tok(ID)"| q226
    q226 -->|"tok(&quot;{&quot;)"| q229
    q227 -.->|"[TokenModeMember]"| q228
    q228 --> q231
    q229 --> q227
    q229 --> q230
    q230 --> q232
    q231 --> q229
    q232 -->|"tok(&quot;}&quot;)"| q233
    q233 --> q31
```

## TokenModeMember

```mermaid
flowchart TD
    q32(["TokenModeMember__Start (32)<br/>RuleStart"])
    q33(["TokenModeMember__Stop (33)<br/>RuleStop"])
    q234["TokenModeMember__Basic_0 (234)<br/>Basic<br/>"]
    q235["TokenModeMember__Basic_1 (235)<br/>Basic<br/>"]
    q236["TokenModeMember__Basic_2 (236)<br/>Basic<br/>"]
    q237["TokenModeMember__Basic_3 (237)<br/>Basic<br/>"]
    q238["TokenModeMember__Basic_4 (238)<br/>Basic<br/>"]
    q239["TokenModeMember__Basic_5 (239)<br/>Basic<br/>"]
    q240["TokenModeMember__Basic_6 (240)<br/>Basic<br/>"]
    q241["TokenModeMember__Basic_7 (241)<br/>Basic<br/>"]
    q242["TokenModeMember__Basic_8 (242)<br/>Basic<br/>"]
    q243["TokenModeMember__Basic_9 (243)<br/>Basic<br/>"]
    q244{"TokenModeMember__Basic_10 (244)<br/>Basic<br/><br/>dec=21"}
    q245["TokenModeMember__BlockEnd (245)<br/>BlockEnd<br/>"]

    q32 --> q244
    q234 -.->|"[TokenDeclUsage]"| q235
    q235 --> q245
    q236 -.->|"[TokenGroupUsage]"| q237
    q237 --> q245
    q238 -.->|"[TokenUsage]"| q239
    q239 --> q245
    q240 -.->|"[KeywordUsage]"| q241
    q241 --> q245
    q242 -.->|"[KeywordSelector]"| q243
    q243 --> q245
    q244 --> q234
    q244 --> q236
    q244 --> q238
    q244 --> q240
    q244 --> q242
    q245 --> q33
```

## TokenDeclUsage

```mermaid
flowchart TD
    q34(["TokenDeclUsage__Start (34)<br/>RuleStart"])
    q35(["TokenDeclUsage__Stop (35)<br/>RuleStop"])
    q246["TokenDeclUsage__Basic_0 (246)<br/>Basic<br/>"]
    q247["TokenDeclUsage__Basic_1 (247)<br/>Basic<br/>"]

    q34 --> q246
    q246 -.->|"[TokenDecl]"| q247
    q247 --> q35
```

## TokenGroupUsage

```mermaid
flowchart TD
    q36(["TokenGroupUsage__Start (36)<br/>RuleStart"])
    q37(["TokenGroupUsage__Stop (37)<br/>RuleStop"])
    q248["TokenGroupUsage__Basic_0 (248)<br/>Basic<br/>"]
    q249["TokenGroupUsage__Basic_1 (249)<br/>Basic<br/>"]

    q36 --> q248
    q248 -.->|"[TokenGroup]"| q249
    q249 --> q37
```

## TokenUsage

```mermaid
flowchart TD
    q38(["TokenUsage__Start (38)<br/>RuleStart"])
    q39(["TokenUsage__Stop (39)<br/>RuleStop"])
    q250["TokenUsage_Modifier_TokenModifier (250)<br/>Basic<br/>"]
    q251["TokenUsage__Basic_0 (251)<br/>Basic<br/>"]
    q252{"TokenUsage__Basic_1 (252)<br/>Basic<br/><br/>dec=22"}
    q253["TokenUsage_TokenRef_ID (253)<br/>Basic<br/>"]
    q254["TokenUsage__Basic_2 (254)<br/>Basic<br/>"]
    q255["TokenUsage__Basic_3 (255)<br/>Basic<br/>"]
    q256{"TokenUsage__Basic_4 (256)<br/>Basic<br/><br/>dec=23"}
    q257["TokenUsage_Semicolon (257)<br/>Basic<br/>"]
    q258["TokenUsage__Basic_5 (258)<br/>Basic<br/>"]
    q259{"TokenUsage__Basic_6 (259)<br/>Basic<br/><br/>dec=24"}

    q38 --> q252
    q250 -->|"tok(TokenModifier)"| q251
    q251 --> q253
    q252 --> q250
    q252 --> q251
    q253 -->|"tok(ID)"| q256
    q254 -.->|"[TokenCommand]"| q255
    q255 --> q259
    q256 --> q254
    q256 --> q255
    q257 -->|"tok(&quot;;&quot;)"| q258
    q258 --> q39
    q259 --> q257
    q259 --> q258
```

## KeywordUsage

```mermaid
flowchart TD
    q40(["KeywordUsage__Start (40)<br/>RuleStart"])
    q41(["KeywordUsage__Stop (41)<br/>RuleStop"])
    q260["KeywordUsage_Modifier_TokenModifier (260)<br/>Basic<br/>"]
    q261["KeywordUsage__Basic_0 (261)<br/>Basic<br/>"]
    q262{"KeywordUsage__Basic_1 (262)<br/>Basic<br/><br/>dec=25"}
    q263["KeywordUsage__Basic_2 (263)<br/>Basic<br/>"]
    q264["KeywordUsage__Basic_3 (264)<br/>Basic<br/>"]
    q265["KeywordUsage__Basic_4 (265)<br/>Basic<br/>"]
    q266{"KeywordUsage__Basic_5 (266)<br/>Basic<br/><br/>dec=26"}
    q267["KeywordUsage_Semicolon (267)<br/>Basic<br/>"]
    q268["KeywordUsage__Basic_6 (268)<br/>Basic<br/>"]
    q269{"KeywordUsage__Basic_7 (269)<br/>Basic<br/><br/>dec=27"}

    q40 --> q262
    q260 -->|"tok(TokenModifier)"| q261
    q261 --> q263
    q262 --> q260
    q262 --> q261
    q263 -.->|"[Keyword]"| q266
    q264 -.->|"[TokenCommand]"| q265
    q265 --> q269
    q266 --> q264
    q266 --> q265
    q267 -->|"tok(&quot;;&quot;)"| q268
    q268 --> q41
    q269 --> q267
    q269 --> q268
```

## KeywordSelector

```mermaid
flowchart TD
    q42(["KeywordSelector__Start (42)<br/>RuleStart"])
    q43(["KeywordSelector__Stop (43)<br/>RuleStop"])
    q270["KeywordSelector_keywords (270)<br/>Basic<br/>"]
    q271["KeywordSelector_Selector_RegexLiteral (271)<br/>Basic<br/>"]
    q272["KeywordSelector_Semicolon (272)<br/>Basic<br/>"]
    q273["KeywordSelector__Basic_0 (273)<br/>Basic<br/>"]
    q274{"KeywordSelector__Basic_1 (274)<br/>Basic<br/><br/>dec=28"}

    q42 --> q270
    q270 -->|"tok(&quot;keywords&quot;)"| q271
    q271 -->|"tok(RegexLiteral)"| q274
    q272 -->|"tok(&quot;;&quot;)"| q273
    q273 --> q43
    q274 --> q272
    q274 --> q273
```

## Alternatives

```mermaid
flowchart TD
    q44(["Alternatives__Start (44)<br/>RuleStart"])
    q45(["Alternatives__Stop (45)<br/>RuleStop"])
    q275["Alternatives__Basic_0 (275)<br/>Basic<br/>"]
    q276["Alternatives_Pipe (276)<br/>Basic<br/>"]
    q277["Alternatives__Basic_1 (277)<br/>Basic<br/>"]
    q278["Alternatives__Basic_2 (278)<br/>Basic<br/>"]
    q279{"Alternatives__LoopBack (279)<br/>LoopBack<br/><br/>dec=29"}
    q280["Alternatives__LoopEnd (280)<br/>LoopEnd<br/>"]
    q281{"Alternatives__Basic_3 (281)<br/>Basic<br/><br/>dec=30"}

    q44 --> q275
    q275 -.->|"[Group]"| q281
    q276 -->|"tok(&quot;|&quot;)"| q277
    q277 -.->|"[Group]"| q278
    q278 --> q279
    q279 --> q276
    q279 --> q280
    q280 --> q45
    q281 --> q276
    q281 --> q280
```

## Group

```mermaid
flowchart TD
    q46(["Group__Start (46)<br/>RuleStart"])
    q47(["Group__Stop (47)<br/>RuleStop"])
    q282["Group__Basic_0 (282)<br/>Basic<br/>"]
    q283["Group__Basic_1 (283)<br/>Basic<br/>"]
    q284["Group__Basic_2 (284)<br/>Basic<br/>"]
    q285{"Group__LoopBack (285)<br/>LoopBack<br/><br/>dec=31"}
    q286["Group__LoopEnd (286)<br/>LoopEnd<br/>"]
    q287{"Group__Basic_3 (287)<br/>Basic<br/><br/>dec=32"}

    q46 --> q282
    q282 -.->|"[Element]"| q287
    q283 -.->|"[Element]"| q284
    q284 --> q285
    q285 --> q283
    q285 --> q286
    q286 --> q47
    q287 --> q283
    q287 --> q286
```

## Element

```mermaid
flowchart TD
    q48(["Element__Start (48)<br/>RuleStart"])
    q49(["Element__Stop (49)<br/>RuleStop"])
    q288["Element__Basic_0 (288)<br/>Basic<br/>"]
    q289["Element__Basic_1 (289)<br/>Basic<br/>"]
    q290["Element__Basic_2 (290)<br/>Basic<br/>"]
    q291["Element__Basic_3 (291)<br/>Basic<br/>"]
    q292["Element__Basic_4 (292)<br/>Basic<br/>"]
    q293["Element__Basic_5 (293)<br/>Basic<br/>"]
    q294["Element__Basic_6 (294)<br/>Basic<br/>"]
    q295["Element__Basic_7 (295)<br/>Basic<br/>"]
    q296["Element_LeftParen (296)<br/>Basic<br/>"]
    q297["Element__Basic_8 (297)<br/>Basic<br/>"]
    q298["Element_RightParen (298)<br/>Basic<br/>"]
    q299["Element__Basic_9 (299)<br/>Basic<br/>"]
    q300{"Element__Basic_10 (300)<br/>Basic<br/><br/>dec=33"}
    q301["Element__BlockEnd (301)<br/>BlockEnd<br/>"]
    q302["Element_Cardinality_Cardinality (302)<br/>Basic<br/>"]
    q303["Element__Basic_11 (303)<br/>Basic<br/>"]
    q304{"Element__Basic_12 (304)<br/>Basic<br/><br/>dec=34"}

    q48 --> q300
    q288 -.->|"[Keyword]"| q289
    q289 --> q301
    q290 -.->|"[Assignment]"| q291
    q291 --> q301
    q292 -.->|"[RuleCall]"| q293
    q293 --> q301
    q294 -.->|"[Action]"| q295
    q295 --> q301
    q296 -->|"tok(&quot;(&quot;)"| q297
    q297 -.->|"[Alternatives]"| q298
    q298 -->|"tok(&quot;)&quot;)"| q299
    q299 --> q301
    q300 --> q288
    q300 --> q290
    q300 --> q292
    q300 --> q294
    q300 --> q296
    q301 --> q304
    q302 -->|"tok(Cardinality)"| q303
    q303 --> q49
    q304 --> q302
    q304 --> q303
```

## Keyword

```mermaid
flowchart TD
    q50(["Keyword__Start (50)<br/>RuleStart"])
    q51(["Keyword__Stop (51)<br/>RuleStop"])
    q305["Keyword_Value_StringLiteral (305)<br/>Basic<br/>"]
    q306["Keyword__Basic (306)<br/>Basic<br/>"]

    q50 --> q305
    q305 -->|"tok(StringLiteral)"| q306
    q306 --> q51
```

## Assignment

```mermaid
flowchart TD
    q52(["Assignment__Start (52)<br/>RuleStart"])
    q53(["Assignment__Stop (53)<br/>RuleStop"])
    q307["Assignment_Property_ID (307)<br/>Basic<br/>"]
    q308["Assignment_Operator_PlusEquals (308)<br/>Basic<br/>"]
    q309["Assignment__Basic_0 (309)<br/>Basic<br/>"]
    q310["Assignment_Operator_Equals (310)<br/>Basic<br/>"]
    q311["Assignment__Basic_1 (311)<br/>Basic<br/>"]
    q312["Assignment_Operator_QuestionEquals (312)<br/>Basic<br/>"]
    q313["Assignment__Basic_2 (313)<br/>Basic<br/>"]
    q314{"Assignment__Basic_3 (314)<br/>Basic<br/><br/>dec=35"}
    q315["Assignment__BlockEnd (315)<br/>BlockEnd<br/>"]
    q316["Assignment__Basic_4 (316)<br/>Basic<br/>"]
    q317["Assignment__Basic_5 (317)<br/>Basic<br/>"]

    q52 --> q307
    q307 -->|"tok(ID)"| q314
    q308 -->|"tok(&quot;+=&quot;)"| q309
    q309 --> q315
    q310 -->|"tok(&quot;=&quot;)"| q311
    q311 --> q315
    q312 -->|"tok(&quot;?=&quot;)"| q313
    q313 --> q315
    q314 --> q308
    q314 --> q310
    q314 --> q312
    q315 --> q316
    q316 -.->|"[Assignable]"| q317
    q317 --> q53
```

## Assignable

```mermaid
flowchart TD
    q54(["Assignable__Start (54)<br/>RuleStart"])
    q55(["Assignable__Stop (55)<br/>RuleStop"])
    q318["Assignable__Basic_0 (318)<br/>Basic<br/>"]
    q319["Assignable__Basic_1 (319)<br/>Basic<br/>"]
    q320["Assignable__Basic_2 (320)<br/>Basic<br/>"]
    q321["Assignable__Basic_3 (321)<br/>Basic<br/>"]
    q322["Assignable__Basic_4 (322)<br/>Basic<br/>"]
    q323["Assignable__Basic_5 (323)<br/>Basic<br/>"]
    q324["Assignable_LeftParen (324)<br/>Basic<br/>"]
    q325["Assignable__Basic_6 (325)<br/>Basic<br/>"]
    q326["Assignable_RightParen (326)<br/>Basic<br/>"]
    q327["Assignable__Basic_7 (327)<br/>Basic<br/>"]
    q328{"Assignable__Basic_8 (328)<br/>Basic<br/><br/>dec=36"}
    q329["Assignable__BlockEnd (329)<br/>BlockEnd<br/>"]

    q54 --> q328
    q318 -.->|"[Keyword]"| q319
    q319 --> q329
    q320 -.->|"[RuleCall]"| q321
    q321 --> q329
    q322 -.->|"[CrossRef]"| q323
    q323 --> q329
    q324 -->|"tok(&quot;(&quot;)"| q325
    q325 -.->|"[AssignableAlternatives]"| q326
    q326 -->|"tok(&quot;)&quot;)"| q327
    q327 --> q329
    q328 --> q318
    q328 --> q320
    q328 --> q322
    q328 --> q324
    q329 --> q55
```

## AssignableWithoutAlts

```mermaid
flowchart TD
    q56(["AssignableWithoutAlts__Start (56)<br/>RuleStart"])
    q57(["AssignableWithoutAlts__Stop (57)<br/>RuleStop"])
    q330["AssignableWithoutAlts__Basic_0 (330)<br/>Basic<br/>"]
    q331["AssignableWithoutAlts__Basic_1 (331)<br/>Basic<br/>"]
    q332["AssignableWithoutAlts__Basic_2 (332)<br/>Basic<br/>"]
    q333["AssignableWithoutAlts__Basic_3 (333)<br/>Basic<br/>"]
    q334["AssignableWithoutAlts__Basic_4 (334)<br/>Basic<br/>"]
    q335["AssignableWithoutAlts__Basic_5 (335)<br/>Basic<br/>"]
    q336{"AssignableWithoutAlts__Basic_6 (336)<br/>Basic<br/><br/>dec=37"}
    q337["AssignableWithoutAlts__BlockEnd (337)<br/>BlockEnd<br/>"]

    q56 --> q336
    q330 -.->|"[Keyword]"| q331
    q331 --> q337
    q332 -.->|"[RuleCall]"| q333
    q333 --> q337
    q334 -.->|"[CrossRef]"| q335
    q335 --> q337
    q336 --> q330
    q336 --> q332
    q336 --> q334
    q337 --> q57
```

## AssignableAlternatives

```mermaid
flowchart TD
    q58(["AssignableAlternatives__Start (58)<br/>RuleStart"])
    q59(["AssignableAlternatives__Stop (59)<br/>RuleStop"])
    q338["AssignableAlternatives__Basic_0 (338)<br/>Basic<br/>"]
    q339["AssignableAlternatives_Pipe (339)<br/>Basic<br/>"]
    q340["AssignableAlternatives__Basic_1 (340)<br/>Basic<br/>"]
    q341["AssignableAlternatives__Basic_2 (341)<br/>Basic<br/>"]
    q342{"AssignableAlternatives__LoopBack (342)<br/>LoopBack<br/><br/>dec=38"}
    q343["AssignableAlternatives__LoopEnd (343)<br/>LoopEnd<br/>"]
    q344{"AssignableAlternatives__Basic_3 (344)<br/>Basic<br/><br/>dec=39"}

    q58 --> q338
    q338 -.->|"[AssignableWithoutAlts]"| q344
    q339 -->|"tok(&quot;|&quot;)"| q340
    q340 -.->|"[AssignableWithoutAlts]"| q341
    q341 --> q342
    q342 --> q339
    q342 --> q343
    q343 --> q59
    q344 --> q339
    q344 --> q343
```

## CrossRef

```mermaid
flowchart TD
    q60(["CrossRef__Start (60)<br/>RuleStart"])
    q61(["CrossRef__Stop (61)<br/>RuleStop"])
    q345["CrossRef_LeftBracket (345)<br/>Basic<br/>"]
    q346["CrossRef_Type_ID (346)<br/>Basic<br/>"]
    q347["CrossRef_Colon (347)<br/>Basic<br/>"]
    q348["CrossRef__Basic_0 (348)<br/>Basic<br/>"]
    q349["CrossRef__Basic_1 (349)<br/>Basic<br/>"]
    q350{"CrossRef__Basic_2 (350)<br/>Basic<br/><br/>dec=40"}
    q351["CrossRef_RightBracket (351)<br/>Basic<br/>"]
    q352["CrossRef__Basic_3 (352)<br/>Basic<br/>"]

    q60 --> q345
    q345 -->|"tok(&quot;[&quot;)"| q346
    q346 -->|"tok(ID)"| q350
    q347 -->|"tok(&quot;:&quot;)"| q348
    q348 -.->|"[RuleCall]"| q349
    q349 --> q351
    q350 --> q347
    q350 --> q349
    q351 -->|"tok(&quot;]&quot;)"| q352
    q352 --> q61
```

## RuleCall

```mermaid
flowchart TD
    q62(["RuleCall__Start (62)<br/>RuleStart"])
    q63(["RuleCall__Stop (63)<br/>RuleStop"])
    q353["RuleCall_Rule_ID (353)<br/>Basic<br/>"]
    q354["RuleCall__Basic (354)<br/>Basic<br/>"]

    q62 --> q353
    q353 -->|"tok(ID)"| q354
    q354 --> q63
```

## Action

```mermaid
flowchart TD
    q64(["Action__Start (64)<br/>RuleStart"])
    q65(["Action__Stop (65)<br/>RuleStop"])
    q355["Action_LeftBrace (355)<br/>Basic<br/>"]
    q356["Action_Type_ID (356)<br/>Basic<br/>"]
    q357["Action_Dot (357)<br/>Basic<br/>"]
    q358["Action_Property_ID (358)<br/>Basic<br/>"]
    q359["Action_Operator_PlusEquals (359)<br/>Basic<br/>"]
    q360["Action__Basic_0 (360)<br/>Basic<br/>"]
    q361["Action_Operator_Equals (361)<br/>Basic<br/>"]
    q362["Action__Basic_1 (362)<br/>Basic<br/>"]
    q363{"Action__Basic_2 (363)<br/>Basic<br/><br/>dec=41"}
    q364["Action__BlockEnd (364)<br/>BlockEnd<br/>"]
    q365["Action_current (365)<br/>Basic<br/>"]
    q366["Action__Basic_3 (366)<br/>Basic<br/>"]
    q367{"Action__Basic_4 (367)<br/>Basic<br/><br/>dec=42"}
    q368["Action_RightBrace (368)<br/>Basic<br/>"]
    q369["Action__Basic_5 (369)<br/>Basic<br/>"]

    q64 --> q355
    q355 -->|"tok(&quot;{&quot;)"| q356
    q356 -->|"tok(ID)"| q367
    q357 -->|"tok(&quot;.&quot;)"| q358
    q358 -->|"tok(ID)"| q363
    q359 -->|"tok(&quot;+=&quot;)"| q360
    q360 --> q364
    q361 -->|"tok(&quot;=&quot;)"| q362
    q362 --> q364
    q363 --> q359
    q363 --> q361
    q364 --> q365
    q365 -->|"tok(&quot;current&quot;)"| q366
    q366 --> q368
    q367 --> q357
    q367 --> q366
    q368 -->|"tok(&quot;}&quot;)"| q369
    q369 --> q65
```

## CompositeRule

```mermaid
flowchart TD
    q66(["CompositeRule__Start (66)<br/>RuleStart"])
    q67(["CompositeRule__Stop (67)<br/>RuleStop"])
    q370["CompositeRule_composite (370)<br/>Basic<br/>"]
    q371["CompositeRule_Name_ID (371)<br/>Basic<br/>"]
    q372["CompositeRule_Colon (372)<br/>Basic<br/>"]
    q373["CompositeRule__Basic_0 (373)<br/>Basic<br/>"]
    q374["CompositeRule_Semicolon (374)<br/>Basic<br/>"]
    q375["CompositeRule__Basic_1 (375)<br/>Basic<br/>"]
    q376{"CompositeRule__Basic_2 (376)<br/>Basic<br/><br/>dec=43"}

    q66 --> q370
    q370 -->|"tok(&quot;composite&quot;)"| q371
    q371 -->|"tok(ID)"| q372
    q372 -->|"tok(&quot;:&quot;)"| q373
    q373 -.->|"[CompositeAlternatives]"| q376
    q374 -->|"tok(&quot;;&quot;)"| q375
    q375 --> q67
    q376 --> q374
    q376 --> q375
```

## CompositeAlternatives

```mermaid
flowchart TD
    q68(["CompositeAlternatives__Start (68)<br/>RuleStart"])
    q69(["CompositeAlternatives__Stop (69)<br/>RuleStop"])
    q377["CompositeAlternatives__Basic_0 (377)<br/>Basic<br/>"]
    q378["CompositeAlternatives_Pipe (378)<br/>Basic<br/>"]
    q379["CompositeAlternatives__Basic_1 (379)<br/>Basic<br/>"]
    q380["CompositeAlternatives__Basic_2 (380)<br/>Basic<br/>"]
    q381{"CompositeAlternatives__LoopBack (381)<br/>LoopBack<br/><br/>dec=44"}
    q382["CompositeAlternatives__LoopEnd (382)<br/>LoopEnd<br/>"]
    q383{"CompositeAlternatives__Basic_3 (383)<br/>Basic<br/><br/>dec=45"}

    q68 --> q377
    q377 -.->|"[CompositeGroup]"| q383
    q378 -->|"tok(&quot;|&quot;)"| q379
    q379 -.->|"[CompositeGroup]"| q380
    q380 --> q381
    q381 --> q378
    q381 --> q382
    q382 --> q69
    q383 --> q378
    q383 --> q382
```

## CompositeGroup

```mermaid
flowchart TD
    q70(["CompositeGroup__Start (70)<br/>RuleStart"])
    q71(["CompositeGroup__Stop (71)<br/>RuleStop"])
    q384["CompositeGroup__Basic_0 (384)<br/>Basic<br/>"]
    q385["CompositeGroup__Basic_1 (385)<br/>Basic<br/>"]
    q386["CompositeGroup__Basic_2 (386)<br/>Basic<br/>"]
    q387{"CompositeGroup__LoopBack (387)<br/>LoopBack<br/><br/>dec=46"}
    q388["CompositeGroup__LoopEnd (388)<br/>LoopEnd<br/>"]
    q389{"CompositeGroup__Basic_3 (389)<br/>Basic<br/><br/>dec=47"}

    q70 --> q384
    q384 -.->|"[CompositeElement]"| q389
    q385 -.->|"[CompositeElement]"| q386
    q386 --> q387
    q387 --> q385
    q387 --> q388
    q388 --> q71
    q389 --> q385
    q389 --> q388
```

## CompositeElement

```mermaid
flowchart TD
    q72(["CompositeElement__Start (72)<br/>RuleStart"])
    q73(["CompositeElement__Stop (73)<br/>RuleStop"])
    q390["CompositeElement__Basic_0 (390)<br/>Basic<br/>"]
    q391["CompositeElement__Basic_1 (391)<br/>Basic<br/>"]
    q392["CompositeElement__Basic_2 (392)<br/>Basic<br/>"]
    q393["CompositeElement__Basic_3 (393)<br/>Basic<br/>"]
    q394["CompositeElement_LeftParen (394)<br/>Basic<br/>"]
    q395["CompositeElement__Basic_4 (395)<br/>Basic<br/>"]
    q396["CompositeElement_RightParen (396)<br/>Basic<br/>"]
    q397["CompositeElement__Basic_5 (397)<br/>Basic<br/>"]
    q398{"CompositeElement__Basic_6 (398)<br/>Basic<br/><br/>dec=48"}
    q399["CompositeElement__BlockEnd (399)<br/>BlockEnd<br/>"]
    q400["CompositeElement_Cardinality_Cardinality (400)<br/>Basic<br/>"]
    q401["CompositeElement__Basic_7 (401)<br/>Basic<br/>"]
    q402{"CompositeElement__Basic_8 (402)<br/>Basic<br/><br/>dec=49"}

    q72 --> q398
    q390 -.->|"[Keyword]"| q391
    q391 --> q399
    q392 -.->|"[RuleCall]"| q393
    q393 --> q399
    q394 -->|"tok(&quot;(&quot;)"| q395
    q395 -.->|"[CompositeAlternatives]"| q396
    q396 -->|"tok(&quot;)&quot;)"| q397
    q397 --> q399
    q398 --> q390
    q398 --> q392
    q398 --> q394
    q399 --> q402
    q400 -->|"tok(Cardinality)"| q401
    q401 --> q73
    q402 --> q400
    q402 --> q401
```

## InfixRule

```mermaid
flowchart TD
    q74(["InfixRule__Start (74)<br/>RuleStart"])
    q75(["InfixRule__Stop (75)<br/>RuleStop"])
    q403["InfixRule_infix (403)<br/>Basic<br/>"]
    q404["InfixRule_Name_ID (404)<br/>Basic<br/>"]
    q405["InfixRule_on (405)<br/>Basic<br/>"]
    q406["InfixRule__Basic_0 (406)<br/>Basic<br/>"]
    q407["InfixRule_returns (407)<br/>Basic<br/>"]
    q408["InfixRule_ReturnType_ID (408)<br/>Basic<br/>"]
    q409["InfixRule__Basic_1 (409)<br/>Basic<br/>"]
    q410{"InfixRule__Basic_2 (410)<br/>Basic<br/><br/>dec=50"}
    q411["InfixRule_Colon (411)<br/>Basic<br/>"]
    q412["InfixRule__Basic_3 (412)<br/>Basic<br/>"]
    q413["InfixRule_GreaterThan (413)<br/>Basic<br/>"]
    q414["InfixRule__Basic_4 (414)<br/>Basic<br/>"]
    q415["InfixRule__Basic_5 (415)<br/>Basic<br/>"]
    q416{"InfixRule__LoopEntry (416)<br/>LoopEntry<br/><br/>dec=51"}
    q417["InfixRule__LoopEnd (417)<br/>LoopEnd<br/>"]
    q418["InfixRule__LoopBack (418)<br/>LoopBack<br/>"]
    q419["InfixRule_Semicolon (419)<br/>Basic<br/>"]
    q420["InfixRule__Basic_6 (420)<br/>Basic<br/>"]
    q421{"InfixRule__Basic_7 (421)<br/>Basic<br/><br/>dec=52"}

    q74 --> q403
    q403 -->|"tok(&quot;infix&quot;)"| q404
    q404 -->|"tok(ID)"| q405
    q405 -->|"tok(&quot;on&quot;)"| q406
    q406 -.->|"[RuleCall]"| q410
    q407 -->|"tok(&quot;returns&quot;)"| q408
    q408 -->|"tok(ID)"| q409
    q409 --> q411
    q410 --> q407
    q410 --> q409
    q411 -->|"tok(&quot;:&quot;)"| q412
    q412 -.->|"[PrecedenceGroup]"| q416
    q413 -->|"tok(&quot;>&quot;)"| q414
    q414 -.->|"[PrecedenceGroup]"| q415
    q415 --> q418
    q416 --> q413
    q416 --> q417
    q417 --> q421
    q418 --> q416
    q419 -->|"tok(&quot;;&quot;)"| q420
    q420 --> q75
    q421 --> q419
    q421 --> q420
```

## PrecedenceGroup

```mermaid
flowchart TD
    q76(["PrecedenceGroup__Start (76)<br/>RuleStart"])
    q77(["PrecedenceGroup__Stop (77)<br/>RuleStop"])
    q422["PrecedenceGroup_Associativity_left (422)<br/>Basic<br/>"]
    q423["PrecedenceGroup__Basic_0 (423)<br/>Basic<br/>"]
    q424["PrecedenceGroup_Associativity_right (424)<br/>Basic<br/>"]
    q425["PrecedenceGroup__Basic_1 (425)<br/>Basic<br/>"]
    q426{"PrecedenceGroup__Basic_2 (426)<br/>Basic<br/><br/>dec=53"}
    q427["PrecedenceGroup__BlockEnd (427)<br/>BlockEnd<br/>"]
    q428{"PrecedenceGroup__Basic_3 (428)<br/>Basic<br/><br/>dec=54"}
    q429["PrecedenceGroup__Basic_4 (429)<br/>Basic<br/>"]
    q430["PrecedenceGroup_Pipe (430)<br/>Basic<br/>"]
    q431["PrecedenceGroup__Basic_5 (431)<br/>Basic<br/>"]
    q432["PrecedenceGroup__Basic_6 (432)<br/>Basic<br/>"]
    q433{"PrecedenceGroup__LoopEntry (433)<br/>LoopEntry<br/><br/>dec=55"}
    q434["PrecedenceGroup__LoopEnd (434)<br/>LoopEnd<br/>"]
    q435["PrecedenceGroup__LoopBack (435)<br/>LoopBack<br/>"]

    q76 --> q428
    q422 -->|"tok(&quot;left&quot;)"| q423
    q423 --> q427
    q424 -->|"tok(&quot;right&quot;)"| q425
    q425 --> q427
    q426 --> q422
    q426 --> q424
    q427 --> q429
    q428 --> q426
    q428 --> q427
    q429 -.->|"[InfixOperator]"| q433
    q430 -->|"tok(&quot;|&quot;)"| q431
    q431 -.->|"[InfixOperator]"| q432
    q432 --> q435
    q433 --> q430
    q433 --> q434
    q434 --> q77
    q435 --> q433
```

## InfixOperator

```mermaid
flowchart TD
    q78(["InfixOperator__Start (78)<br/>RuleStart"])
    q79(["InfixOperator__Stop (79)<br/>RuleStop"])
    q436["InfixOperator__Basic_0 (436)<br/>Basic<br/>"]
    q437["InfixOperator__Basic_1 (437)<br/>Basic<br/>"]
    q438["InfixOperator__Basic_2 (438)<br/>Basic<br/>"]
    q439["InfixOperator__Basic_3 (439)<br/>Basic<br/>"]
    q440{"InfixOperator__Basic_4 (440)<br/>Basic<br/><br/>dec=56"}
    q441["InfixOperator__BlockEnd (441)<br/>BlockEnd<br/>"]

    q78 --> q440
    q436 -.->|"[Keyword]"| q437
    q437 --> q441
    q438 -.->|"[RuleCall]"| q439
    q439 --> q441
    q440 --> q436
    q440 --> q438
    q441 --> q79
```

