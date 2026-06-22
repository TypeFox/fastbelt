# Runtime ATN for grammar

## Grammar

```mermaid
flowchart TD
    q0(["Grammar__Start (0)<br/>RuleStart"])
    q1(["Grammar__Stop (1)<br/>RuleStop"])
    q52["Grammar_grammar (52)<br/>Basic<br/>"]
    q53["Grammar_Name_ID (53)<br/>Basic<br/>"]
    q54["Grammar_Semicolon (54)<br/>Basic<br/>"]
    q55["Grammar__Basic_0 (55)<br/>Basic<br/>"]
    q56["Grammar__Basic_1 (56)<br/>Basic<br/>"]
    q57["Grammar__Basic_2 (57)<br/>Basic<br/>"]
    q58["Grammar__Basic_3 (58)<br/>Basic<br/>"]
    q59["Grammar__Basic_4 (59)<br/>Basic<br/>"]
    q60["Grammar__Basic_5 (60)<br/>Basic<br/>"]
    q61["Grammar__Basic_6 (61)<br/>Basic<br/>"]
    q62["Grammar__Basic_7 (62)<br/>Basic<br/>"]
    q63["Grammar__Basic_8 (63)<br/>Basic<br/>"]
    q64["Grammar__Basic_9 (64)<br/>Basic<br/>"]
    q65{"Grammar__Basic_10 (65)<br/>Basic<br/><br/>dec=0"}
    q66["Grammar__BlockEnd (66)<br/>BlockEnd<br/>"]
    q67{"Grammar__LoopEntry (67)<br/>LoopEntry<br/><br/>dec=1"}
    q68["Grammar__LoopEnd (68)<br/>LoopEnd<br/>"]
    q69["Grammar__LoopBack (69)<br/>LoopBack<br/>"]

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
    q2(["Interface__Start (2)<br/>RuleStart"])
    q3(["Interface__Stop (3)<br/>RuleStop"])
    q70["Interface_interface (70)<br/>Basic<br/>"]
    q71["Interface_Name_ID (71)<br/>Basic<br/>"]
    q72["Interface_extends (72)<br/>Basic<br/>"]
    q73["Interface_Extends_ID_0 (73)<br/>Basic<br/>"]
    q74["Interface_Comma (74)<br/>Basic<br/>"]
    q75["Interface_Extends_ID_1 (75)<br/>Basic<br/>"]
    q76["Interface__Basic_0 (76)<br/>Basic<br/>"]
    q77{"Interface__LoopEntry_0 (77)<br/>LoopEntry<br/><br/>dec=2"}
    q78["Interface__LoopEnd_0 (78)<br/>LoopEnd<br/>"]
    q79["Interface__LoopBack_0 (79)<br/>LoopBack<br/>"]
    q80{"Interface__Basic_1 (80)<br/>Basic<br/><br/>dec=3"}
    q81["Interface_LeftBrace (81)<br/>Basic<br/>"]
    q82["Interface__Basic_2 (82)<br/>Basic<br/>"]
    q83["Interface__Basic_3 (83)<br/>Basic<br/>"]
    q84{"Interface__LoopEntry_1 (84)<br/>LoopEntry<br/><br/>dec=4"}
    q85["Interface__LoopEnd_1 (85)<br/>LoopEnd<br/>"]
    q86["Interface__LoopBack_1 (86)<br/>LoopBack<br/>"]
    q87["Interface_RightBrace (87)<br/>Basic<br/>"]
    q88["Interface__Basic_4 (88)<br/>Basic<br/>"]

    q2 --> q70
    q70 -->|"tok(&quot;interface&quot;)"| q71
    q71 -->|"tok(ID)"| q80
    q72 -->|"tok(&quot;extends&quot;)"| q73
    q73 -->|"tok(ID)"| q77
    q74 -->|"tok(&quot;,&quot;)"| q75
    q75 -->|"tok(ID)"| q76
    q76 --> q79
    q77 --> q74
    q77 --> q78
    q78 --> q81
    q79 --> q77
    q80 --> q72
    q80 --> q78
    q81 -->|"tok(&quot;{&quot;)"| q84
    q82 -.->|"[Field]"| q83
    q83 --> q86
    q84 --> q82
    q84 --> q85
    q85 --> q87
    q86 --> q84
    q87 -->|"tok(&quot;}&quot;)"| q88
    q88 --> q3
```

## Field

```mermaid
flowchart TD
    q4(["Field__Start (4)<br/>RuleStart"])
    q5(["Field__Stop (5)<br/>RuleStop"])
    q89["Field_Name_ID (89)<br/>Basic<br/>"]
    q90["Field__Basic_0 (90)<br/>Basic<br/>"]
    q91["Field__Basic_1 (91)<br/>Basic<br/>"]

    q4 --> q89
    q89 -->|"tok(ID)"| q90
    q90 -.->|"[FieldType]"| q91
    q91 --> q5
```

## FieldType

```mermaid
flowchart TD
    q6(["FieldType__Start (6)<br/>RuleStart"])
    q7(["FieldType__Stop (7)<br/>RuleStop"])
    q92["FieldType__Basic_0 (92)<br/>Basic<br/>"]
    q93["FieldType__Basic_1 (93)<br/>Basic<br/>"]
    q94["FieldType__Basic_2 (94)<br/>Basic<br/>"]
    q95["FieldType__Basic_3 (95)<br/>Basic<br/>"]
    q96["FieldType__Basic_4 (96)<br/>Basic<br/>"]
    q97["FieldType__Basic_5 (97)<br/>Basic<br/>"]
    q98["FieldType__Basic_6 (98)<br/>Basic<br/>"]
    q99["FieldType__Basic_7 (99)<br/>Basic<br/>"]
    q100{"FieldType__Basic_8 (100)<br/>Basic<br/><br/>dec=5"}
    q101["FieldType__BlockEnd (101)<br/>BlockEnd<br/>"]

    q6 --> q100
    q92 -.->|"[SimpleType]"| q93
    q93 --> q101
    q94 -.->|"[ReferenceType]"| q95
    q95 --> q101
    q96 -.->|"[ArrayType]"| q97
    q97 --> q101
    q98 -.->|"[PrimitiveType]"| q99
    q99 --> q101
    q100 --> q92
    q100 --> q94
    q100 --> q96
    q100 --> q98
    q101 --> q7
```

## ArrayType

```mermaid
flowchart TD
    q8(["ArrayType__Start (8)<br/>RuleStart"])
    q9(["ArrayType__Stop (9)<br/>RuleStop"])
    q102["ArrayType_LeftBracket (102)<br/>Basic<br/>"]
    q103["ArrayType_RightBracket (103)<br/>Basic<br/>"]
    q104["ArrayType__Basic_0 (104)<br/>Basic<br/>"]
    q105["ArrayType__Basic_1 (105)<br/>Basic<br/>"]

    q8 --> q102
    q102 -->|"tok(&quot;[&quot;)"| q103
    q103 -->|"tok(&quot;]&quot;)"| q104
    q104 -.->|"[FieldType]"| q105
    q105 --> q9
```

## ReferenceType

```mermaid
flowchart TD
    q10(["ReferenceType__Start (10)<br/>RuleStart"])
    q11(["ReferenceType__Stop (11)<br/>RuleStop"])
    q106["ReferenceType_Asterisk (106)<br/>Basic<br/>"]
    q107["ReferenceType_Type_ID (107)<br/>Basic<br/>"]
    q108["ReferenceType__Basic (108)<br/>Basic<br/>"]

    q10 --> q106
    q106 -->|"tok(&quot;*&quot;)"| q107
    q107 -->|"tok(ID)"| q108
    q108 --> q11
```

## SimpleType

```mermaid
flowchart TD
    q12(["SimpleType__Start (12)<br/>RuleStart"])
    q13(["SimpleType__Stop (13)<br/>RuleStop"])
    q109["SimpleType_Type_ID (109)<br/>Basic<br/>"]
    q110["SimpleType__Basic (110)<br/>Basic<br/>"]

    q12 --> q109
    q109 -->|"tok(ID)"| q110
    q110 --> q13
```

## PrimitiveType

```mermaid
flowchart TD
    q14(["PrimitiveType__Start (14)<br/>RuleStart"])
    q15(["PrimitiveType__Stop (15)<br/>RuleStop"])
    q111["PrimitiveType_Type_string (111)<br/>Basic<br/>"]
    q112["PrimitiveType__Basic_0 (112)<br/>Basic<br/>"]
    q113["PrimitiveType_Type_bool (113)<br/>Basic<br/>"]
    q114["PrimitiveType__Basic_1 (114)<br/>Basic<br/>"]
    q115["PrimitiveType_Type_composite (115)<br/>Basic<br/>"]
    q116["PrimitiveType__Basic_2 (116)<br/>Basic<br/>"]
    q117{"PrimitiveType__Basic_3 (117)<br/>Basic<br/><br/>dec=6"}
    q118["PrimitiveType__BlockEnd (118)<br/>BlockEnd<br/>"]

    q14 --> q117
    q111 -->|"tok(&quot;string&quot;)"| q112
    q112 --> q118
    q113 -->|"tok(&quot;bool&quot;)"| q114
    q114 --> q118
    q115 -->|"tok(&quot;composite&quot;)"| q116
    q116 --> q118
    q117 --> q111
    q117 --> q113
    q117 --> q115
    q118 --> q15
```

## ParserRule

```mermaid
flowchart TD
    q16(["ParserRule__Start (16)<br/>RuleStart"])
    q17(["ParserRule__Stop (17)<br/>RuleStop"])
    q119["ParserRule_Name_ID (119)<br/>Basic<br/>"]
    q120["ParserRule_returns (120)<br/>Basic<br/>"]
    q121["ParserRule_ReturnType_ID (121)<br/>Basic<br/>"]
    q122["ParserRule__Basic_0 (122)<br/>Basic<br/>"]
    q123{"ParserRule__Basic_1 (123)<br/>Basic<br/><br/>dec=7"}
    q124["ParserRule_Colon (124)<br/>Basic<br/>"]
    q125["ParserRule__Basic_2 (125)<br/>Basic<br/>"]
    q126["ParserRule_Semicolon (126)<br/>Basic<br/>"]
    q127["ParserRule__Basic_3 (127)<br/>Basic<br/>"]

    q16 --> q119
    q119 -->|"tok(ID)"| q123
    q120 -->|"tok(&quot;returns&quot;)"| q121
    q121 -->|"tok(ID)"| q122
    q122 --> q124
    q123 --> q120
    q123 --> q122
    q124 -->|"tok(&quot;:&quot;)"| q125
    q125 -.->|"[Alternatives]"| q126
    q126 -->|"tok(&quot;;&quot;)"| q127
    q127 --> q17
```

## Token

```mermaid
flowchart TD
    q18(["Token__Start (18)<br/>RuleStart"])
    q19(["Token__Stop (19)<br/>RuleStop"])
    q128["Token_Type_hidden (128)<br/>Basic<br/>"]
    q129["Token__Basic_0 (129)<br/>Basic<br/>"]
    q130["Token_Type_comment (130)<br/>Basic<br/>"]
    q131["Token__Basic_1 (131)<br/>Basic<br/>"]
    q132{"Token__Basic_2 (132)<br/>Basic<br/><br/>dec=8"}
    q133["Token__BlockEnd (133)<br/>BlockEnd<br/>"]
    q134{"Token__Basic_3 (134)<br/>Basic<br/><br/>dec=9"}
    q135["Token_token (135)<br/>Basic<br/>"]
    q136["Token_Name_ID (136)<br/>Basic<br/>"]
    q137["Token_Colon (137)<br/>Basic<br/>"]
    q138["Token_Regexp_RegexLiteral (138)<br/>Basic<br/>"]
    q139["Token_Semicolon (139)<br/>Basic<br/>"]
    q140["Token__Basic_4 (140)<br/>Basic<br/>"]

    q18 --> q134
    q128 -->|"tok(&quot;hidden&quot;)"| q129
    q129 --> q133
    q130 -->|"tok(&quot;comment&quot;)"| q131
    q131 --> q133
    q132 --> q128
    q132 --> q130
    q133 --> q135
    q134 --> q132
    q134 --> q133
    q135 -->|"tok(&quot;token&quot;)"| q136
    q136 -->|"tok(ID)"| q137
    q137 -->|"tok(&quot;:&quot;)"| q138
    q138 -->|"tok(RegexLiteral)"| q139
    q139 -->|"tok(&quot;;&quot;)"| q140
    q140 --> q19
```

## TokenGroup

```mermaid
flowchart TD
    q20(["TokenGroup__Start (20)<br/>RuleStart"])
    q21(["TokenGroup__Stop (21)<br/>RuleStop"])
    q141["TokenGroup_token (141)<br/>Basic<br/>"]
    q142["TokenGroup_group (142)<br/>Basic<br/>"]
    q143["TokenGroup_Name_ID (143)<br/>Basic<br/>"]
    q144["TokenGroup_LeftBrace (144)<br/>Basic<br/>"]
    q145["TokenGroup_TokenRefs_ID (145)<br/>Basic<br/>"]
    q146["TokenGroup__Basic_0 (146)<br/>Basic<br/>"]
    q147["TokenGroup_keywords (147)<br/>Basic<br/>"]
    q148["TokenGroup_Regexps_RegexLiteral (148)<br/>Basic<br/>"]
    q149["TokenGroup__Basic_1 (149)<br/>Basic<br/>"]
    q150["TokenGroup__Basic_2 (150)<br/>Basic<br/>"]
    q151["TokenGroup__Basic_3 (151)<br/>Basic<br/>"]
    q152{"TokenGroup__Basic_4 (152)<br/>Basic<br/><br/>dec=10"}
    q153["TokenGroup__BlockEnd (153)<br/>BlockEnd<br/>"]
    q154{"TokenGroup__LoopEntry (154)<br/>LoopEntry<br/><br/>dec=11"}
    q155["TokenGroup__LoopEnd (155)<br/>LoopEnd<br/>"]
    q156["TokenGroup__LoopBack (156)<br/>LoopBack<br/>"]
    q157["TokenGroup_RightBrace (157)<br/>Basic<br/>"]
    q158["TokenGroup__Basic_5 (158)<br/>Basic<br/>"]

    q20 --> q141
    q141 -->|"tok(&quot;token&quot;)"| q142
    q142 -->|"tok(&quot;group&quot;)"| q143
    q143 -->|"tok(ID)"| q144
    q144 -->|"tok(&quot;{&quot;)"| q154
    q145 -->|"tok(ID)"| q146
    q146 --> q153
    q147 -->|"tok(&quot;keywords&quot;)"| q148
    q148 -->|"tok(RegexLiteral)"| q149
    q149 --> q153
    q150 -.->|"[Keyword]"| q151
    q151 --> q153
    q152 --> q145
    q152 --> q147
    q152 --> q150
    q153 --> q156
    q154 --> q152
    q154 --> q155
    q155 --> q157
    q156 --> q154
    q157 -->|"tok(&quot;}&quot;)"| q158
    q158 --> q21
```

## Alternatives

```mermaid
flowchart TD
    q22(["Alternatives__Start (22)<br/>RuleStart"])
    q23(["Alternatives__Stop (23)<br/>RuleStop"])
    q159["Alternatives__Basic_0 (159)<br/>Basic<br/>"]
    q160["Alternatives_Pipe (160)<br/>Basic<br/>"]
    q161["Alternatives__Basic_1 (161)<br/>Basic<br/>"]
    q162["Alternatives__Basic_2 (162)<br/>Basic<br/>"]
    q163{"Alternatives__LoopBack (163)<br/>LoopBack<br/><br/>dec=12"}
    q164["Alternatives__LoopEnd (164)<br/>LoopEnd<br/>"]
    q165{"Alternatives__Basic_3 (165)<br/>Basic<br/><br/>dec=13"}

    q22 --> q159
    q159 -.->|"[Group]"| q165
    q160 -->|"tok(&quot;|&quot;)"| q161
    q161 -.->|"[Group]"| q162
    q162 --> q163
    q163 --> q160
    q163 --> q164
    q164 --> q23
    q165 --> q160
    q165 --> q164
```

## Group

```mermaid
flowchart TD
    q24(["Group__Start (24)<br/>RuleStart"])
    q25(["Group__Stop (25)<br/>RuleStop"])
    q166["Group__Basic_0 (166)<br/>Basic<br/>"]
    q167["Group__Basic_1 (167)<br/>Basic<br/>"]
    q168["Group__Basic_2 (168)<br/>Basic<br/>"]
    q169{"Group__LoopBack (169)<br/>LoopBack<br/><br/>dec=14"}
    q170["Group__LoopEnd (170)<br/>LoopEnd<br/>"]
    q171{"Group__Basic_3 (171)<br/>Basic<br/><br/>dec=15"}

    q24 --> q166
    q166 -.->|"[Element]"| q171
    q167 -.->|"[Element]"| q168
    q168 --> q169
    q169 --> q167
    q169 --> q170
    q170 --> q25
    q171 --> q167
    q171 --> q170
```

## Element

```mermaid
flowchart TD
    q26(["Element__Start (26)<br/>RuleStart"])
    q27(["Element__Stop (27)<br/>RuleStop"])
    q172["Element__Basic_0 (172)<br/>Basic<br/>"]
    q173["Element__Basic_1 (173)<br/>Basic<br/>"]
    q174["Element__Basic_2 (174)<br/>Basic<br/>"]
    q175["Element__Basic_3 (175)<br/>Basic<br/>"]
    q176["Element__Basic_4 (176)<br/>Basic<br/>"]
    q177["Element__Basic_5 (177)<br/>Basic<br/>"]
    q178["Element__Basic_6 (178)<br/>Basic<br/>"]
    q179["Element__Basic_7 (179)<br/>Basic<br/>"]
    q180["Element_LeftParen (180)<br/>Basic<br/>"]
    q181["Element__Basic_8 (181)<br/>Basic<br/>"]
    q182["Element_RightParen (182)<br/>Basic<br/>"]
    q183["Element__Basic_9 (183)<br/>Basic<br/>"]
    q184{"Element__Basic_10 (184)<br/>Basic<br/><br/>dec=16"}
    q185["Element__BlockEnd_0 (185)<br/>BlockEnd<br/>"]
    q186["Element_Cardinality_Asterisk (186)<br/>Basic<br/>"]
    q187["Element__Basic_11 (187)<br/>Basic<br/>"]
    q188["Element_Cardinality_Plus (188)<br/>Basic<br/>"]
    q189["Element__Basic_12 (189)<br/>Basic<br/>"]
    q190["Element_Cardinality_Question (190)<br/>Basic<br/>"]
    q191["Element__Basic_13 (191)<br/>Basic<br/>"]
    q192{"Element__Basic_14 (192)<br/>Basic<br/><br/>dec=17"}
    q193["Element__BlockEnd_1 (193)<br/>BlockEnd<br/>"]
    q194{"Element__Basic_15 (194)<br/>Basic<br/><br/>dec=18"}

    q26 --> q184
    q172 -.->|"[Keyword]"| q173
    q173 --> q185
    q174 -.->|"[Assignment]"| q175
    q175 --> q185
    q176 -.->|"[RuleCall]"| q177
    q177 --> q185
    q178 -.->|"[Action]"| q179
    q179 --> q185
    q180 -->|"tok(&quot;(&quot;)"| q181
    q181 -.->|"[Alternatives]"| q182
    q182 -->|"tok(&quot;)&quot;)"| q183
    q183 --> q185
    q184 --> q172
    q184 --> q174
    q184 --> q176
    q184 --> q178
    q184 --> q180
    q185 --> q194
    q186 -->|"tok(&quot;*&quot;)"| q187
    q187 --> q193
    q188 -->|"tok(&quot;+&quot;)"| q189
    q189 --> q193
    q190 -->|"tok(&quot;?&quot;)"| q191
    q191 --> q193
    q192 --> q186
    q192 --> q188
    q192 --> q190
    q193 --> q27
    q194 --> q192
    q194 --> q193
```

## Keyword

```mermaid
flowchart TD
    q28(["Keyword__Start (28)<br/>RuleStart"])
    q29(["Keyword__Stop (29)<br/>RuleStop"])
    q195["Keyword_Value_StringLiteral (195)<br/>Basic<br/>"]
    q196["Keyword__Basic (196)<br/>Basic<br/>"]

    q28 --> q195
    q195 -->|"tok(StringLiteral)"| q196
    q196 --> q29
```

## Assignment

```mermaid
flowchart TD
    q30(["Assignment__Start (30)<br/>RuleStart"])
    q31(["Assignment__Stop (31)<br/>RuleStop"])
    q197["Assignment_Property_ID (197)<br/>Basic<br/>"]
    q198["Assignment_Operator_PlusEquals (198)<br/>Basic<br/>"]
    q199["Assignment__Basic_0 (199)<br/>Basic<br/>"]
    q200["Assignment_Operator_Equals (200)<br/>Basic<br/>"]
    q201["Assignment__Basic_1 (201)<br/>Basic<br/>"]
    q202["Assignment_Operator_QuestionEquals (202)<br/>Basic<br/>"]
    q203["Assignment__Basic_2 (203)<br/>Basic<br/>"]
    q204{"Assignment__Basic_3 (204)<br/>Basic<br/><br/>dec=19"}
    q205["Assignment__BlockEnd (205)<br/>BlockEnd<br/>"]
    q206["Assignment__Basic_4 (206)<br/>Basic<br/>"]
    q207["Assignment__Basic_5 (207)<br/>Basic<br/>"]

    q30 --> q197
    q197 -->|"tok(ID)"| q204
    q198 -->|"tok(&quot;+=&quot;)"| q199
    q199 --> q205
    q200 -->|"tok(&quot;=&quot;)"| q201
    q201 --> q205
    q202 -->|"tok(&quot;?=&quot;)"| q203
    q203 --> q205
    q204 --> q198
    q204 --> q200
    q204 --> q202
    q205 --> q206
    q206 -.->|"[Assignable]"| q207
    q207 --> q31
```

## Assignable

```mermaid
flowchart TD
    q32(["Assignable__Start (32)<br/>RuleStart"])
    q33(["Assignable__Stop (33)<br/>RuleStop"])
    q208["Assignable__Basic_0 (208)<br/>Basic<br/>"]
    q209["Assignable__Basic_1 (209)<br/>Basic<br/>"]
    q210["Assignable__Basic_2 (210)<br/>Basic<br/>"]
    q211["Assignable__Basic_3 (211)<br/>Basic<br/>"]
    q212["Assignable__Basic_4 (212)<br/>Basic<br/>"]
    q213["Assignable__Basic_5 (213)<br/>Basic<br/>"]
    q214["Assignable_LeftParen (214)<br/>Basic<br/>"]
    q215["Assignable__Basic_6 (215)<br/>Basic<br/>"]
    q216["Assignable_RightParen (216)<br/>Basic<br/>"]
    q217["Assignable__Basic_7 (217)<br/>Basic<br/>"]
    q218{"Assignable__Basic_8 (218)<br/>Basic<br/><br/>dec=20"}
    q219["Assignable__BlockEnd (219)<br/>BlockEnd<br/>"]

    q32 --> q218
    q208 -.->|"[Keyword]"| q209
    q209 --> q219
    q210 -.->|"[RuleCall]"| q211
    q211 --> q219
    q212 -.->|"[CrossRef]"| q213
    q213 --> q219
    q214 -->|"tok(&quot;(&quot;)"| q215
    q215 -.->|"[AssignableAlternatives]"| q216
    q216 -->|"tok(&quot;)&quot;)"| q217
    q217 --> q219
    q218 --> q208
    q218 --> q210
    q218 --> q212
    q218 --> q214
    q219 --> q33
```

## AssignableWithoutAlts

```mermaid
flowchart TD
    q34(["AssignableWithoutAlts__Start (34)<br/>RuleStart"])
    q35(["AssignableWithoutAlts__Stop (35)<br/>RuleStop"])
    q220["AssignableWithoutAlts__Basic_0 (220)<br/>Basic<br/>"]
    q221["AssignableWithoutAlts__Basic_1 (221)<br/>Basic<br/>"]
    q222["AssignableWithoutAlts__Basic_2 (222)<br/>Basic<br/>"]
    q223["AssignableWithoutAlts__Basic_3 (223)<br/>Basic<br/>"]
    q224["AssignableWithoutAlts__Basic_4 (224)<br/>Basic<br/>"]
    q225["AssignableWithoutAlts__Basic_5 (225)<br/>Basic<br/>"]
    q226{"AssignableWithoutAlts__Basic_6 (226)<br/>Basic<br/><br/>dec=21"}
    q227["AssignableWithoutAlts__BlockEnd (227)<br/>BlockEnd<br/>"]

    q34 --> q226
    q220 -.->|"[Keyword]"| q221
    q221 --> q227
    q222 -.->|"[RuleCall]"| q223
    q223 --> q227
    q224 -.->|"[CrossRef]"| q225
    q225 --> q227
    q226 --> q220
    q226 --> q222
    q226 --> q224
    q227 --> q35
```

## AssignableAlternatives

```mermaid
flowchart TD
    q36(["AssignableAlternatives__Start (36)<br/>RuleStart"])
    q37(["AssignableAlternatives__Stop (37)<br/>RuleStop"])
    q228["AssignableAlternatives__Basic_0 (228)<br/>Basic<br/>"]
    q229["AssignableAlternatives_Pipe (229)<br/>Basic<br/>"]
    q230["AssignableAlternatives__Basic_1 (230)<br/>Basic<br/>"]
    q231["AssignableAlternatives__Basic_2 (231)<br/>Basic<br/>"]
    q232{"AssignableAlternatives__LoopBack (232)<br/>LoopBack<br/><br/>dec=22"}
    q233["AssignableAlternatives__LoopEnd (233)<br/>LoopEnd<br/>"]
    q234{"AssignableAlternatives__Basic_3 (234)<br/>Basic<br/><br/>dec=23"}

    q36 --> q228
    q228 -.->|"[AssignableWithoutAlts]"| q234
    q229 -->|"tok(&quot;|&quot;)"| q230
    q230 -.->|"[AssignableWithoutAlts]"| q231
    q231 --> q232
    q232 --> q229
    q232 --> q233
    q233 --> q37
    q234 --> q229
    q234 --> q233
```

## CrossRef

```mermaid
flowchart TD
    q38(["CrossRef__Start (38)<br/>RuleStart"])
    q39(["CrossRef__Stop (39)<br/>RuleStop"])
    q235["CrossRef_LeftBracket (235)<br/>Basic<br/>"]
    q236["CrossRef_Type_ID (236)<br/>Basic<br/>"]
    q237["CrossRef_Colon (237)<br/>Basic<br/>"]
    q238["CrossRef__Basic_0 (238)<br/>Basic<br/>"]
    q239["CrossRef__Basic_1 (239)<br/>Basic<br/>"]
    q240{"CrossRef__Basic_2 (240)<br/>Basic<br/><br/>dec=24"}
    q241["CrossRef_RightBracket (241)<br/>Basic<br/>"]
    q242["CrossRef__Basic_3 (242)<br/>Basic<br/>"]

    q38 --> q235
    q235 -->|"tok(&quot;[&quot;)"| q236
    q236 -->|"tok(ID)"| q240
    q237 -->|"tok(&quot;:&quot;)"| q238
    q238 -.->|"[RuleCall]"| q239
    q239 --> q241
    q240 --> q237
    q240 --> q239
    q241 -->|"tok(&quot;]&quot;)"| q242
    q242 --> q39
```

## RuleCall

```mermaid
flowchart TD
    q40(["RuleCall__Start (40)<br/>RuleStart"])
    q41(["RuleCall__Stop (41)<br/>RuleStop"])
    q243["RuleCall_Rule_ID (243)<br/>Basic<br/>"]
    q244["RuleCall__Basic (244)<br/>Basic<br/>"]

    q40 --> q243
    q243 -->|"tok(ID)"| q244
    q244 --> q41
```

## Action

```mermaid
flowchart TD
    q42(["Action__Start (42)<br/>RuleStart"])
    q43(["Action__Stop (43)<br/>RuleStop"])
    q245["Action_LeftBrace (245)<br/>Basic<br/>"]
    q246["Action_Type_ID (246)<br/>Basic<br/>"]
    q247["Action_Dot (247)<br/>Basic<br/>"]
    q248["Action_Property_ID (248)<br/>Basic<br/>"]
    q249["Action_Operator_PlusEquals (249)<br/>Basic<br/>"]
    q250["Action__Basic_0 (250)<br/>Basic<br/>"]
    q251["Action_Operator_Equals (251)<br/>Basic<br/>"]
    q252["Action__Basic_1 (252)<br/>Basic<br/>"]
    q253{"Action__Basic_2 (253)<br/>Basic<br/><br/>dec=25"}
    q254["Action__BlockEnd (254)<br/>BlockEnd<br/>"]
    q255["Action_current (255)<br/>Basic<br/>"]
    q256["Action__Basic_3 (256)<br/>Basic<br/>"]
    q257{"Action__Basic_4 (257)<br/>Basic<br/><br/>dec=26"}
    q258["Action_RightBrace (258)<br/>Basic<br/>"]
    q259["Action__Basic_5 (259)<br/>Basic<br/>"]

    q42 --> q245
    q245 -->|"tok(&quot;{&quot;)"| q246
    q246 -->|"tok(ID)"| q257
    q247 -->|"tok(&quot;.&quot;)"| q248
    q248 -->|"tok(ID)"| q253
    q249 -->|"tok(&quot;+=&quot;)"| q250
    q250 --> q254
    q251 -->|"tok(&quot;=&quot;)"| q252
    q252 --> q254
    q253 --> q249
    q253 --> q251
    q254 --> q255
    q255 -->|"tok(&quot;current&quot;)"| q256
    q256 --> q258
    q257 --> q247
    q257 --> q256
    q258 -->|"tok(&quot;}&quot;)"| q259
    q259 --> q43
```

## CompositeRule

```mermaid
flowchart TD
    q44(["CompositeRule__Start (44)<br/>RuleStart"])
    q45(["CompositeRule__Stop (45)<br/>RuleStop"])
    q260["CompositeRule_composite (260)<br/>Basic<br/>"]
    q261["CompositeRule_Name_ID (261)<br/>Basic<br/>"]
    q262["CompositeRule_Colon (262)<br/>Basic<br/>"]
    q263["CompositeRule__Basic_0 (263)<br/>Basic<br/>"]
    q264["CompositeRule_Semicolon (264)<br/>Basic<br/>"]
    q265["CompositeRule__Basic_1 (265)<br/>Basic<br/>"]

    q44 --> q260
    q260 -->|"tok(&quot;composite&quot;)"| q261
    q261 -->|"tok(ID)"| q262
    q262 -->|"tok(&quot;:&quot;)"| q263
    q263 -.->|"[CompositeAlternatives]"| q264
    q264 -->|"tok(&quot;;&quot;)"| q265
    q265 --> q45
```

## CompositeAlternatives

```mermaid
flowchart TD
    q46(["CompositeAlternatives__Start (46)<br/>RuleStart"])
    q47(["CompositeAlternatives__Stop (47)<br/>RuleStop"])
    q266["CompositeAlternatives__Basic_0 (266)<br/>Basic<br/>"]
    q267["CompositeAlternatives_Pipe (267)<br/>Basic<br/>"]
    q268["CompositeAlternatives__Basic_1 (268)<br/>Basic<br/>"]
    q269["CompositeAlternatives__Basic_2 (269)<br/>Basic<br/>"]
    q270{"CompositeAlternatives__LoopBack (270)<br/>LoopBack<br/><br/>dec=27"}
    q271["CompositeAlternatives__LoopEnd (271)<br/>LoopEnd<br/>"]
    q272{"CompositeAlternatives__Basic_3 (272)<br/>Basic<br/><br/>dec=28"}

    q46 --> q266
    q266 -.->|"[CompositeGroup]"| q272
    q267 -->|"tok(&quot;|&quot;)"| q268
    q268 -.->|"[CompositeGroup]"| q269
    q269 --> q270
    q270 --> q267
    q270 --> q271
    q271 --> q47
    q272 --> q267
    q272 --> q271
```

## CompositeGroup

```mermaid
flowchart TD
    q48(["CompositeGroup__Start (48)<br/>RuleStart"])
    q49(["CompositeGroup__Stop (49)<br/>RuleStop"])
    q273["CompositeGroup__Basic_0 (273)<br/>Basic<br/>"]
    q274["CompositeGroup__Basic_1 (274)<br/>Basic<br/>"]
    q275["CompositeGroup__Basic_2 (275)<br/>Basic<br/>"]
    q276{"CompositeGroup__LoopBack (276)<br/>LoopBack<br/><br/>dec=29"}
    q277["CompositeGroup__LoopEnd (277)<br/>LoopEnd<br/>"]
    q278{"CompositeGroup__Basic_3 (278)<br/>Basic<br/><br/>dec=30"}

    q48 --> q273
    q273 -.->|"[CompositeElement]"| q278
    q274 -.->|"[CompositeElement]"| q275
    q275 --> q276
    q276 --> q274
    q276 --> q277
    q277 --> q49
    q278 --> q274
    q278 --> q277
```

## CompositeElement

```mermaid
flowchart TD
    q50(["CompositeElement__Start (50)<br/>RuleStart"])
    q51(["CompositeElement__Stop (51)<br/>RuleStop"])
    q279["CompositeElement__Basic_0 (279)<br/>Basic<br/>"]
    q280["CompositeElement__Basic_1 (280)<br/>Basic<br/>"]
    q281["CompositeElement__Basic_2 (281)<br/>Basic<br/>"]
    q282["CompositeElement__Basic_3 (282)<br/>Basic<br/>"]
    q283["CompositeElement_LeftParen (283)<br/>Basic<br/>"]
    q284["CompositeElement__Basic_4 (284)<br/>Basic<br/>"]
    q285["CompositeElement_RightParen (285)<br/>Basic<br/>"]
    q286["CompositeElement__Basic_5 (286)<br/>Basic<br/>"]
    q287{"CompositeElement__Basic_6 (287)<br/>Basic<br/><br/>dec=31"}
    q288["CompositeElement__BlockEnd_0 (288)<br/>BlockEnd<br/>"]
    q289["CompositeElement_Cardinality_Asterisk (289)<br/>Basic<br/>"]
    q290["CompositeElement__Basic_7 (290)<br/>Basic<br/>"]
    q291["CompositeElement_Cardinality_Plus (291)<br/>Basic<br/>"]
    q292["CompositeElement__Basic_8 (292)<br/>Basic<br/>"]
    q293["CompositeElement_Cardinality_Question (293)<br/>Basic<br/>"]
    q294["CompositeElement__Basic_9 (294)<br/>Basic<br/>"]
    q295{"CompositeElement__Basic_10 (295)<br/>Basic<br/><br/>dec=32"}
    q296["CompositeElement__BlockEnd_1 (296)<br/>BlockEnd<br/>"]
    q297{"CompositeElement__Basic_11 (297)<br/>Basic<br/><br/>dec=33"}

    q50 --> q287
    q279 -.->|"[Keyword]"| q280
    q280 --> q288
    q281 -.->|"[RuleCall]"| q282
    q282 --> q288
    q283 -->|"tok(&quot;(&quot;)"| q284
    q284 -.->|"[CompositeAlternatives]"| q285
    q285 -->|"tok(&quot;)&quot;)"| q286
    q286 --> q288
    q287 --> q279
    q287 --> q281
    q287 --> q283
    q288 --> q297
    q289 -->|"tok(&quot;*&quot;)"| q290
    q290 --> q296
    q291 -->|"tok(&quot;+&quot;)"| q292
    q292 --> q296
    q293 -->|"tok(&quot;?&quot;)"| q294
    q294 --> q296
    q295 --> q289
    q295 --> q291
    q295 --> q293
    q296 --> q51
    q297 --> q295
    q297 --> q296
```

