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
    q56{"Grammar__Basic_1 (56)<br/>Basic<br/><br/>dec=0"}
    q57["Grammar__Basic_2 (57)<br/>Basic<br/>"]
    q58["Grammar__Basic_3 (58)<br/>Basic<br/>"]
    q59["Grammar__Basic_4 (59)<br/>Basic<br/>"]
    q60["Grammar__Basic_5 (60)<br/>Basic<br/>"]
    q61["Grammar__Basic_6 (61)<br/>Basic<br/>"]
    q62["Grammar__Basic_7 (62)<br/>Basic<br/>"]
    q63["Grammar__Basic_8 (63)<br/>Basic<br/>"]
    q64["Grammar__Basic_9 (64)<br/>Basic<br/>"]
    q65["Grammar__Basic_10 (65)<br/>Basic<br/>"]
    q66["Grammar__Basic_11 (66)<br/>Basic<br/>"]
    q67{"Grammar__Basic_12 (67)<br/>Basic<br/><br/>dec=1"}
    q68["Grammar__BlockEnd (68)<br/>BlockEnd<br/>"]
    q69{"Grammar__LoopEntry (69)<br/>LoopEntry<br/><br/>dec=2"}
    q70["Grammar__LoopEnd (70)<br/>LoopEnd<br/>"]
    q71["Grammar__LoopBack (71)<br/>LoopBack<br/>"]

    q0 --> q52
    q52 -->|"tok(&quot;grammar&quot;)"| q53
    q53 -->|"tok(ID)"| q56
    q54 -->|"tok(&quot;;&quot;)"| q55
    q55 --> q69
    q56 --> q54
    q56 --> q55
    q57 -.->|"[ParserRule]"| q58
    q58 --> q68
    q59 -.->|"[Token]"| q60
    q60 --> q68
    q61 -.->|"[TokenGroup]"| q62
    q62 --> q68
    q63 -.->|"[Interface]"| q64
    q64 --> q68
    q65 -.->|"[CompositeRule]"| q66
    q66 --> q68
    q67 --> q57
    q67 --> q59
    q67 --> q61
    q67 --> q63
    q67 --> q65
    q68 --> q71
    q69 --> q67
    q69 --> q70
    q70 --> q1
    q71 --> q69
```

## Interface

```mermaid
flowchart TD
    q2(["Interface__Start (2)<br/>RuleStart"])
    q3(["Interface__Stop (3)<br/>RuleStop"])
    q72["Interface_interface (72)<br/>Basic<br/>"]
    q73["Interface_Name_ID (73)<br/>Basic<br/>"]
    q74["Interface_extends (74)<br/>Basic<br/>"]
    q75["Interface_Extends_ID_0 (75)<br/>Basic<br/>"]
    q76["Interface_Comma (76)<br/>Basic<br/>"]
    q77["Interface_Extends_ID_1 (77)<br/>Basic<br/>"]
    q78["Interface__Basic_0 (78)<br/>Basic<br/>"]
    q79{"Interface__LoopEntry_0 (79)<br/>LoopEntry<br/><br/>dec=3"}
    q80["Interface__LoopEnd_0 (80)<br/>LoopEnd<br/>"]
    q81["Interface__LoopBack_0 (81)<br/>LoopBack<br/>"]
    q82{"Interface__Basic_1 (82)<br/>Basic<br/><br/>dec=4"}
    q83["Interface_LeftBrace (83)<br/>Basic<br/>"]
    q84["Interface__Basic_2 (84)<br/>Basic<br/>"]
    q85["Interface__Basic_3 (85)<br/>Basic<br/>"]
    q86{"Interface__LoopEntry_1 (86)<br/>LoopEntry<br/><br/>dec=5"}
    q87["Interface__LoopEnd_1 (87)<br/>LoopEnd<br/>"]
    q88["Interface__LoopBack_1 (88)<br/>LoopBack<br/>"]
    q89["Interface_RightBrace (89)<br/>Basic<br/>"]
    q90["Interface__Basic_4 (90)<br/>Basic<br/>"]

    q2 --> q72
    q72 -->|"tok(&quot;interface&quot;)"| q73
    q73 -->|"tok(ID)"| q82
    q74 -->|"tok(&quot;extends&quot;)"| q75
    q75 -->|"tok(ID)"| q79
    q76 -->|"tok(&quot;,&quot;)"| q77
    q77 -->|"tok(ID)"| q78
    q78 --> q81
    q79 --> q76
    q79 --> q80
    q80 --> q83
    q81 --> q79
    q82 --> q74
    q82 --> q80
    q83 -->|"tok(&quot;{&quot;)"| q86
    q84 -.->|"[Field]"| q85
    q85 --> q88
    q86 --> q84
    q86 --> q87
    q87 --> q89
    q88 --> q86
    q89 -->|"tok(&quot;}&quot;)"| q90
    q90 --> q3
```

## Field

```mermaid
flowchart TD
    q4(["Field__Start (4)<br/>RuleStart"])
    q5(["Field__Stop (5)<br/>RuleStop"])
    q91["Field_Name_ID (91)<br/>Basic<br/>"]
    q92["Field__Basic_0 (92)<br/>Basic<br/>"]
    q93["Field__Basic_1 (93)<br/>Basic<br/>"]

    q4 --> q91
    q91 -->|"tok(ID)"| q92
    q92 -.->|"[FieldType]"| q93
    q93 --> q5
```

## FieldType

```mermaid
flowchart TD
    q6(["FieldType__Start (6)<br/>RuleStart"])
    q7(["FieldType__Stop (7)<br/>RuleStop"])
    q94["FieldType__Basic_0 (94)<br/>Basic<br/>"]
    q95["FieldType__Basic_1 (95)<br/>Basic<br/>"]
    q96["FieldType__Basic_2 (96)<br/>Basic<br/>"]
    q97["FieldType__Basic_3 (97)<br/>Basic<br/>"]
    q98["FieldType__Basic_4 (98)<br/>Basic<br/>"]
    q99["FieldType__Basic_5 (99)<br/>Basic<br/>"]
    q100["FieldType__Basic_6 (100)<br/>Basic<br/>"]
    q101["FieldType__Basic_7 (101)<br/>Basic<br/>"]
    q102{"FieldType__Basic_8 (102)<br/>Basic<br/><br/>dec=6"}
    q103["FieldType__BlockEnd (103)<br/>BlockEnd<br/>"]

    q6 --> q102
    q94 -.->|"[SimpleType]"| q95
    q95 --> q103
    q96 -.->|"[ReferenceType]"| q97
    q97 --> q103
    q98 -.->|"[ArrayType]"| q99
    q99 --> q103
    q100 -.->|"[PrimitiveType]"| q101
    q101 --> q103
    q102 --> q94
    q102 --> q96
    q102 --> q98
    q102 --> q100
    q103 --> q7
```

## ArrayType

```mermaid
flowchart TD
    q8(["ArrayType__Start (8)<br/>RuleStart"])
    q9(["ArrayType__Stop (9)<br/>RuleStop"])
    q104["ArrayType_LeftBracket (104)<br/>Basic<br/>"]
    q105["ArrayType_RightBracket (105)<br/>Basic<br/>"]
    q106["ArrayType__Basic_0 (106)<br/>Basic<br/>"]
    q107["ArrayType__Basic_1 (107)<br/>Basic<br/>"]

    q8 --> q104
    q104 -->|"tok(&quot;[&quot;)"| q105
    q105 -->|"tok(&quot;]&quot;)"| q106
    q106 -.->|"[FieldType]"| q107
    q107 --> q9
```

## ReferenceType

```mermaid
flowchart TD
    q10(["ReferenceType__Start (10)<br/>RuleStart"])
    q11(["ReferenceType__Stop (11)<br/>RuleStop"])
    q108["ReferenceType_Asterisk (108)<br/>Basic<br/>"]
    q109["ReferenceType_Type_ID (109)<br/>Basic<br/>"]
    q110["ReferenceType__Basic (110)<br/>Basic<br/>"]

    q10 --> q108
    q108 -->|"tok(&quot;*&quot;)"| q109
    q109 -->|"tok(ID)"| q110
    q110 --> q11
```

## SimpleType

```mermaid
flowchart TD
    q12(["SimpleType__Start (12)<br/>RuleStart"])
    q13(["SimpleType__Stop (13)<br/>RuleStop"])
    q111["SimpleType_Type_ID (111)<br/>Basic<br/>"]
    q112["SimpleType__Basic (112)<br/>Basic<br/>"]

    q12 --> q111
    q111 -->|"tok(ID)"| q112
    q112 --> q13
```

## PrimitiveType

```mermaid
flowchart TD
    q14(["PrimitiveType__Start (14)<br/>RuleStart"])
    q15(["PrimitiveType__Stop (15)<br/>RuleStop"])
    q113["PrimitiveType_Type_string (113)<br/>Basic<br/>"]
    q114["PrimitiveType__Basic_0 (114)<br/>Basic<br/>"]
    q115["PrimitiveType_Type_bool (115)<br/>Basic<br/>"]
    q116["PrimitiveType__Basic_1 (116)<br/>Basic<br/>"]
    q117["PrimitiveType_Type_composite (117)<br/>Basic<br/>"]
    q118["PrimitiveType__Basic_2 (118)<br/>Basic<br/>"]
    q119{"PrimitiveType__Basic_3 (119)<br/>Basic<br/><br/>dec=7"}
    q120["PrimitiveType__BlockEnd (120)<br/>BlockEnd<br/>"]

    q14 --> q119
    q113 -->|"tok(&quot;string&quot;)"| q114
    q114 --> q120
    q115 -->|"tok(&quot;bool&quot;)"| q116
    q116 --> q120
    q117 -->|"tok(&quot;composite&quot;)"| q118
    q118 --> q120
    q119 --> q113
    q119 --> q115
    q119 --> q117
    q120 --> q15
```

## ParserRule

```mermaid
flowchart TD
    q16(["ParserRule__Start (16)<br/>RuleStart"])
    q17(["ParserRule__Stop (17)<br/>RuleStop"])
    q121["ParserRule_Name_ID (121)<br/>Basic<br/>"]
    q122["ParserRule_returns (122)<br/>Basic<br/>"]
    q123["ParserRule_ReturnType_ID (123)<br/>Basic<br/>"]
    q124["ParserRule__Basic_0 (124)<br/>Basic<br/>"]
    q125{"ParserRule__Basic_1 (125)<br/>Basic<br/><br/>dec=8"}
    q126["ParserRule_Colon (126)<br/>Basic<br/>"]
    q127["ParserRule__Basic_2 (127)<br/>Basic<br/>"]
    q128["ParserRule_Semicolon (128)<br/>Basic<br/>"]
    q129["ParserRule__Basic_3 (129)<br/>Basic<br/>"]
    q130{"ParserRule__Basic_4 (130)<br/>Basic<br/><br/>dec=9"}

    q16 --> q121
    q121 -->|"tok(ID)"| q125
    q122 -->|"tok(&quot;returns&quot;)"| q123
    q123 -->|"tok(ID)"| q124
    q124 --> q126
    q125 --> q122
    q125 --> q124
    q126 -->|"tok(&quot;:&quot;)"| q127
    q127 -.->|"[Alternatives]"| q130
    q128 -->|"tok(&quot;;&quot;)"| q129
    q129 --> q17
    q130 --> q128
    q130 --> q129
```

## Token

```mermaid
flowchart TD
    q18(["Token__Start (18)<br/>RuleStart"])
    q19(["Token__Stop (19)<br/>RuleStop"])
    q131["Token_Type_hidden (131)<br/>Basic<br/>"]
    q132["Token__Basic_0 (132)<br/>Basic<br/>"]
    q133["Token_Type_comment (133)<br/>Basic<br/>"]
    q134["Token__Basic_1 (134)<br/>Basic<br/>"]
    q135{"Token__Basic_2 (135)<br/>Basic<br/><br/>dec=10"}
    q136["Token__BlockEnd (136)<br/>BlockEnd<br/>"]
    q137{"Token__Basic_3 (137)<br/>Basic<br/><br/>dec=11"}
    q138["Token_token (138)<br/>Basic<br/>"]
    q139["Token_Name_ID (139)<br/>Basic<br/>"]
    q140["Token_Colon (140)<br/>Basic<br/>"]
    q141["Token_Regexp_RegexLiteral (141)<br/>Basic<br/>"]
    q142["Token_Semicolon (142)<br/>Basic<br/>"]
    q143["Token__Basic_4 (143)<br/>Basic<br/>"]
    q144{"Token__Basic_5 (144)<br/>Basic<br/><br/>dec=12"}

    q18 --> q137
    q131 -->|"tok(&quot;hidden&quot;)"| q132
    q132 --> q136
    q133 -->|"tok(&quot;comment&quot;)"| q134
    q134 --> q136
    q135 --> q131
    q135 --> q133
    q136 --> q138
    q137 --> q135
    q137 --> q136
    q138 -->|"tok(&quot;token&quot;)"| q139
    q139 -->|"tok(ID)"| q140
    q140 -->|"tok(&quot;:&quot;)"| q141
    q141 -->|"tok(RegexLiteral)"| q144
    q142 -->|"tok(&quot;;&quot;)"| q143
    q143 --> q19
    q144 --> q142
    q144 --> q143
```

## TokenGroup

```mermaid
flowchart TD
    q20(["TokenGroup__Start (20)<br/>RuleStart"])
    q21(["TokenGroup__Stop (21)<br/>RuleStop"])
    q145["TokenGroup_token (145)<br/>Basic<br/>"]
    q146["TokenGroup_group (146)<br/>Basic<br/>"]
    q147["TokenGroup_Name_ID (147)<br/>Basic<br/>"]
    q148["TokenGroup_LeftBrace (148)<br/>Basic<br/>"]
    q149["TokenGroup_TokenRefs_ID (149)<br/>Basic<br/>"]
    q150["TokenGroup__Basic_0 (150)<br/>Basic<br/>"]
    q151["TokenGroup_keywords (151)<br/>Basic<br/>"]
    q152["TokenGroup_Regexps_RegexLiteral (152)<br/>Basic<br/>"]
    q153["TokenGroup__Basic_1 (153)<br/>Basic<br/>"]
    q154["TokenGroup__Basic_2 (154)<br/>Basic<br/>"]
    q155["TokenGroup__Basic_3 (155)<br/>Basic<br/>"]
    q156{"TokenGroup__Basic_4 (156)<br/>Basic<br/><br/>dec=13"}
    q157["TokenGroup__BlockEnd (157)<br/>BlockEnd<br/>"]
    q158{"TokenGroup__LoopEntry (158)<br/>LoopEntry<br/><br/>dec=14"}
    q159["TokenGroup__LoopEnd (159)<br/>LoopEnd<br/>"]
    q160["TokenGroup__LoopBack (160)<br/>LoopBack<br/>"]
    q161["TokenGroup_RightBrace (161)<br/>Basic<br/>"]
    q162["TokenGroup__Basic_5 (162)<br/>Basic<br/>"]

    q20 --> q145
    q145 -->|"tok(&quot;token&quot;)"| q146
    q146 -->|"tok(&quot;group&quot;)"| q147
    q147 -->|"tok(ID)"| q148
    q148 -->|"tok(&quot;{&quot;)"| q158
    q149 -->|"tok(ID)"| q150
    q150 --> q157
    q151 -->|"tok(&quot;keywords&quot;)"| q152
    q152 -->|"tok(RegexLiteral)"| q153
    q153 --> q157
    q154 -.->|"[Keyword]"| q155
    q155 --> q157
    q156 --> q149
    q156 --> q151
    q156 --> q154
    q157 --> q160
    q158 --> q156
    q158 --> q159
    q159 --> q161
    q160 --> q158
    q161 -->|"tok(&quot;}&quot;)"| q162
    q162 --> q21
```

## Alternatives

```mermaid
flowchart TD
    q22(["Alternatives__Start (22)<br/>RuleStart"])
    q23(["Alternatives__Stop (23)<br/>RuleStop"])
    q163["Alternatives__Basic_0 (163)<br/>Basic<br/>"]
    q164["Alternatives_Pipe (164)<br/>Basic<br/>"]
    q165["Alternatives__Basic_1 (165)<br/>Basic<br/>"]
    q166["Alternatives__Basic_2 (166)<br/>Basic<br/>"]
    q167{"Alternatives__LoopBack (167)<br/>LoopBack<br/><br/>dec=15"}
    q168["Alternatives__LoopEnd (168)<br/>LoopEnd<br/>"]
    q169{"Alternatives__Basic_3 (169)<br/>Basic<br/><br/>dec=16"}

    q22 --> q163
    q163 -.->|"[Group]"| q169
    q164 -->|"tok(&quot;|&quot;)"| q165
    q165 -.->|"[Group]"| q166
    q166 --> q167
    q167 --> q164
    q167 --> q168
    q168 --> q23
    q169 --> q164
    q169 --> q168
```

## Group

```mermaid
flowchart TD
    q24(["Group__Start (24)<br/>RuleStart"])
    q25(["Group__Stop (25)<br/>RuleStop"])
    q170["Group__Basic_0 (170)<br/>Basic<br/>"]
    q171["Group__Basic_1 (171)<br/>Basic<br/>"]
    q172["Group__Basic_2 (172)<br/>Basic<br/>"]
    q173{"Group__LoopBack (173)<br/>LoopBack<br/><br/>dec=17"}
    q174["Group__LoopEnd (174)<br/>LoopEnd<br/>"]
    q175{"Group__Basic_3 (175)<br/>Basic<br/><br/>dec=18"}

    q24 --> q170
    q170 -.->|"[Element]"| q175
    q171 -.->|"[Element]"| q172
    q172 --> q173
    q173 --> q171
    q173 --> q174
    q174 --> q25
    q175 --> q171
    q175 --> q174
```

## Element

```mermaid
flowchart TD
    q26(["Element__Start (26)<br/>RuleStart"])
    q27(["Element__Stop (27)<br/>RuleStop"])
    q176["Element__Basic_0 (176)<br/>Basic<br/>"]
    q177["Element__Basic_1 (177)<br/>Basic<br/>"]
    q178["Element__Basic_2 (178)<br/>Basic<br/>"]
    q179["Element__Basic_3 (179)<br/>Basic<br/>"]
    q180["Element__Basic_4 (180)<br/>Basic<br/>"]
    q181["Element__Basic_5 (181)<br/>Basic<br/>"]
    q182["Element__Basic_6 (182)<br/>Basic<br/>"]
    q183["Element__Basic_7 (183)<br/>Basic<br/>"]
    q184["Element_LeftParen (184)<br/>Basic<br/>"]
    q185["Element__Basic_8 (185)<br/>Basic<br/>"]
    q186["Element_RightParen (186)<br/>Basic<br/>"]
    q187["Element__Basic_9 (187)<br/>Basic<br/>"]
    q188{"Element__Basic_10 (188)<br/>Basic<br/><br/>dec=19"}
    q189["Element__BlockEnd_0 (189)<br/>BlockEnd<br/>"]
    q190["Element_Cardinality_Asterisk (190)<br/>Basic<br/>"]
    q191["Element__Basic_11 (191)<br/>Basic<br/>"]
    q192["Element_Cardinality_Plus (192)<br/>Basic<br/>"]
    q193["Element__Basic_12 (193)<br/>Basic<br/>"]
    q194["Element_Cardinality_Question (194)<br/>Basic<br/>"]
    q195["Element__Basic_13 (195)<br/>Basic<br/>"]
    q196{"Element__Basic_14 (196)<br/>Basic<br/><br/>dec=20"}
    q197["Element__BlockEnd_1 (197)<br/>BlockEnd<br/>"]
    q198{"Element__Basic_15 (198)<br/>Basic<br/><br/>dec=21"}

    q26 --> q188
    q176 -.->|"[Keyword]"| q177
    q177 --> q189
    q178 -.->|"[Assignment]"| q179
    q179 --> q189
    q180 -.->|"[RuleCall]"| q181
    q181 --> q189
    q182 -.->|"[Action]"| q183
    q183 --> q189
    q184 -->|"tok(&quot;(&quot;)"| q185
    q185 -.->|"[Alternatives]"| q186
    q186 -->|"tok(&quot;)&quot;)"| q187
    q187 --> q189
    q188 --> q176
    q188 --> q178
    q188 --> q180
    q188 --> q182
    q188 --> q184
    q189 --> q198
    q190 -->|"tok(&quot;*&quot;)"| q191
    q191 --> q197
    q192 -->|"tok(&quot;+&quot;)"| q193
    q193 --> q197
    q194 -->|"tok(&quot;?&quot;)"| q195
    q195 --> q197
    q196 --> q190
    q196 --> q192
    q196 --> q194
    q197 --> q27
    q198 --> q196
    q198 --> q197
```

## Keyword

```mermaid
flowchart TD
    q28(["Keyword__Start (28)<br/>RuleStart"])
    q29(["Keyword__Stop (29)<br/>RuleStop"])
    q199["Keyword_Value_StringLiteral (199)<br/>Basic<br/>"]
    q200["Keyword__Basic (200)<br/>Basic<br/>"]

    q28 --> q199
    q199 -->|"tok(StringLiteral)"| q200
    q200 --> q29
```

## Assignment

```mermaid
flowchart TD
    q30(["Assignment__Start (30)<br/>RuleStart"])
    q31(["Assignment__Stop (31)<br/>RuleStop"])
    q201["Assignment_Property_ID (201)<br/>Basic<br/>"]
    q202["Assignment_Operator_PlusEquals (202)<br/>Basic<br/>"]
    q203["Assignment__Basic_0 (203)<br/>Basic<br/>"]
    q204["Assignment_Operator_Equals (204)<br/>Basic<br/>"]
    q205["Assignment__Basic_1 (205)<br/>Basic<br/>"]
    q206["Assignment_Operator_QuestionEquals (206)<br/>Basic<br/>"]
    q207["Assignment__Basic_2 (207)<br/>Basic<br/>"]
    q208{"Assignment__Basic_3 (208)<br/>Basic<br/><br/>dec=22"}
    q209["Assignment__BlockEnd (209)<br/>BlockEnd<br/>"]
    q210["Assignment__Basic_4 (210)<br/>Basic<br/>"]
    q211["Assignment__Basic_5 (211)<br/>Basic<br/>"]

    q30 --> q201
    q201 -->|"tok(ID)"| q208
    q202 -->|"tok(&quot;+=&quot;)"| q203
    q203 --> q209
    q204 -->|"tok(&quot;=&quot;)"| q205
    q205 --> q209
    q206 -->|"tok(&quot;?=&quot;)"| q207
    q207 --> q209
    q208 --> q202
    q208 --> q204
    q208 --> q206
    q209 --> q210
    q210 -.->|"[Assignable]"| q211
    q211 --> q31
```

## Assignable

```mermaid
flowchart TD
    q32(["Assignable__Start (32)<br/>RuleStart"])
    q33(["Assignable__Stop (33)<br/>RuleStop"])
    q212["Assignable__Basic_0 (212)<br/>Basic<br/>"]
    q213["Assignable__Basic_1 (213)<br/>Basic<br/>"]
    q214["Assignable__Basic_2 (214)<br/>Basic<br/>"]
    q215["Assignable__Basic_3 (215)<br/>Basic<br/>"]
    q216["Assignable__Basic_4 (216)<br/>Basic<br/>"]
    q217["Assignable__Basic_5 (217)<br/>Basic<br/>"]
    q218["Assignable_LeftParen (218)<br/>Basic<br/>"]
    q219["Assignable__Basic_6 (219)<br/>Basic<br/>"]
    q220["Assignable_RightParen (220)<br/>Basic<br/>"]
    q221["Assignable__Basic_7 (221)<br/>Basic<br/>"]
    q222{"Assignable__Basic_8 (222)<br/>Basic<br/><br/>dec=23"}
    q223["Assignable__BlockEnd (223)<br/>BlockEnd<br/>"]

    q32 --> q222
    q212 -.->|"[Keyword]"| q213
    q213 --> q223
    q214 -.->|"[RuleCall]"| q215
    q215 --> q223
    q216 -.->|"[CrossRef]"| q217
    q217 --> q223
    q218 -->|"tok(&quot;(&quot;)"| q219
    q219 -.->|"[AssignableAlternatives]"| q220
    q220 -->|"tok(&quot;)&quot;)"| q221
    q221 --> q223
    q222 --> q212
    q222 --> q214
    q222 --> q216
    q222 --> q218
    q223 --> q33
```

## AssignableWithoutAlts

```mermaid
flowchart TD
    q34(["AssignableWithoutAlts__Start (34)<br/>RuleStart"])
    q35(["AssignableWithoutAlts__Stop (35)<br/>RuleStop"])
    q224["AssignableWithoutAlts__Basic_0 (224)<br/>Basic<br/>"]
    q225["AssignableWithoutAlts__Basic_1 (225)<br/>Basic<br/>"]
    q226["AssignableWithoutAlts__Basic_2 (226)<br/>Basic<br/>"]
    q227["AssignableWithoutAlts__Basic_3 (227)<br/>Basic<br/>"]
    q228["AssignableWithoutAlts__Basic_4 (228)<br/>Basic<br/>"]
    q229["AssignableWithoutAlts__Basic_5 (229)<br/>Basic<br/>"]
    q230{"AssignableWithoutAlts__Basic_6 (230)<br/>Basic<br/><br/>dec=24"}
    q231["AssignableWithoutAlts__BlockEnd (231)<br/>BlockEnd<br/>"]

    q34 --> q230
    q224 -.->|"[Keyword]"| q225
    q225 --> q231
    q226 -.->|"[RuleCall]"| q227
    q227 --> q231
    q228 -.->|"[CrossRef]"| q229
    q229 --> q231
    q230 --> q224
    q230 --> q226
    q230 --> q228
    q231 --> q35
```

## AssignableAlternatives

```mermaid
flowchart TD
    q36(["AssignableAlternatives__Start (36)<br/>RuleStart"])
    q37(["AssignableAlternatives__Stop (37)<br/>RuleStop"])
    q232["AssignableAlternatives__Basic_0 (232)<br/>Basic<br/>"]
    q233["AssignableAlternatives_Pipe (233)<br/>Basic<br/>"]
    q234["AssignableAlternatives__Basic_1 (234)<br/>Basic<br/>"]
    q235["AssignableAlternatives__Basic_2 (235)<br/>Basic<br/>"]
    q236{"AssignableAlternatives__LoopBack (236)<br/>LoopBack<br/><br/>dec=25"}
    q237["AssignableAlternatives__LoopEnd (237)<br/>LoopEnd<br/>"]
    q238{"AssignableAlternatives__Basic_3 (238)<br/>Basic<br/><br/>dec=26"}

    q36 --> q232
    q232 -.->|"[AssignableWithoutAlts]"| q238
    q233 -->|"tok(&quot;|&quot;)"| q234
    q234 -.->|"[AssignableWithoutAlts]"| q235
    q235 --> q236
    q236 --> q233
    q236 --> q237
    q237 --> q37
    q238 --> q233
    q238 --> q237
```

## CrossRef

```mermaid
flowchart TD
    q38(["CrossRef__Start (38)<br/>RuleStart"])
    q39(["CrossRef__Stop (39)<br/>RuleStop"])
    q239["CrossRef_LeftBracket (239)<br/>Basic<br/>"]
    q240["CrossRef_Type_ID (240)<br/>Basic<br/>"]
    q241["CrossRef_Colon (241)<br/>Basic<br/>"]
    q242["CrossRef__Basic_0 (242)<br/>Basic<br/>"]
    q243["CrossRef__Basic_1 (243)<br/>Basic<br/>"]
    q244{"CrossRef__Basic_2 (244)<br/>Basic<br/><br/>dec=27"}
    q245["CrossRef_RightBracket (245)<br/>Basic<br/>"]
    q246["CrossRef__Basic_3 (246)<br/>Basic<br/>"]

    q38 --> q239
    q239 -->|"tok(&quot;[&quot;)"| q240
    q240 -->|"tok(ID)"| q244
    q241 -->|"tok(&quot;:&quot;)"| q242
    q242 -.->|"[RuleCall]"| q243
    q243 --> q245
    q244 --> q241
    q244 --> q243
    q245 -->|"tok(&quot;]&quot;)"| q246
    q246 --> q39
```

## RuleCall

```mermaid
flowchart TD
    q40(["RuleCall__Start (40)<br/>RuleStart"])
    q41(["RuleCall__Stop (41)<br/>RuleStop"])
    q247["RuleCall_Rule_ID (247)<br/>Basic<br/>"]
    q248["RuleCall__Basic (248)<br/>Basic<br/>"]

    q40 --> q247
    q247 -->|"tok(ID)"| q248
    q248 --> q41
```

## Action

```mermaid
flowchart TD
    q42(["Action__Start (42)<br/>RuleStart"])
    q43(["Action__Stop (43)<br/>RuleStop"])
    q249["Action_LeftBrace (249)<br/>Basic<br/>"]
    q250["Action_Type_ID (250)<br/>Basic<br/>"]
    q251["Action_Dot (251)<br/>Basic<br/>"]
    q252["Action_Property_ID (252)<br/>Basic<br/>"]
    q253["Action_Operator_PlusEquals (253)<br/>Basic<br/>"]
    q254["Action__Basic_0 (254)<br/>Basic<br/>"]
    q255["Action_Operator_Equals (255)<br/>Basic<br/>"]
    q256["Action__Basic_1 (256)<br/>Basic<br/>"]
    q257{"Action__Basic_2 (257)<br/>Basic<br/><br/>dec=28"}
    q258["Action__BlockEnd (258)<br/>BlockEnd<br/>"]
    q259["Action_current (259)<br/>Basic<br/>"]
    q260["Action__Basic_3 (260)<br/>Basic<br/>"]
    q261{"Action__Basic_4 (261)<br/>Basic<br/><br/>dec=29"}
    q262["Action_RightBrace (262)<br/>Basic<br/>"]
    q263["Action__Basic_5 (263)<br/>Basic<br/>"]

    q42 --> q249
    q249 -->|"tok(&quot;{&quot;)"| q250
    q250 -->|"tok(ID)"| q261
    q251 -->|"tok(&quot;.&quot;)"| q252
    q252 -->|"tok(ID)"| q257
    q253 -->|"tok(&quot;+=&quot;)"| q254
    q254 --> q258
    q255 -->|"tok(&quot;=&quot;)"| q256
    q256 --> q258
    q257 --> q253
    q257 --> q255
    q258 --> q259
    q259 -->|"tok(&quot;current&quot;)"| q260
    q260 --> q262
    q261 --> q251
    q261 --> q260
    q262 -->|"tok(&quot;}&quot;)"| q263
    q263 --> q43
```

## CompositeRule

```mermaid
flowchart TD
    q44(["CompositeRule__Start (44)<br/>RuleStart"])
    q45(["CompositeRule__Stop (45)<br/>RuleStop"])
    q264["CompositeRule_composite (264)<br/>Basic<br/>"]
    q265["CompositeRule_Name_ID (265)<br/>Basic<br/>"]
    q266["CompositeRule_Colon (266)<br/>Basic<br/>"]
    q267["CompositeRule__Basic_0 (267)<br/>Basic<br/>"]
    q268["CompositeRule_Semicolon (268)<br/>Basic<br/>"]
    q269["CompositeRule__Basic_1 (269)<br/>Basic<br/>"]
    q270{"CompositeRule__Basic_2 (270)<br/>Basic<br/><br/>dec=30"}

    q44 --> q264
    q264 -->|"tok(&quot;composite&quot;)"| q265
    q265 -->|"tok(ID)"| q266
    q266 -->|"tok(&quot;:&quot;)"| q267
    q267 -.->|"[CompositeAlternatives]"| q270
    q268 -->|"tok(&quot;;&quot;)"| q269
    q269 --> q45
    q270 --> q268
    q270 --> q269
```

## CompositeAlternatives

```mermaid
flowchart TD
    q46(["CompositeAlternatives__Start (46)<br/>RuleStart"])
    q47(["CompositeAlternatives__Stop (47)<br/>RuleStop"])
    q271["CompositeAlternatives__Basic_0 (271)<br/>Basic<br/>"]
    q272["CompositeAlternatives_Pipe (272)<br/>Basic<br/>"]
    q273["CompositeAlternatives__Basic_1 (273)<br/>Basic<br/>"]
    q274["CompositeAlternatives__Basic_2 (274)<br/>Basic<br/>"]
    q275{"CompositeAlternatives__LoopBack (275)<br/>LoopBack<br/><br/>dec=31"}
    q276["CompositeAlternatives__LoopEnd (276)<br/>LoopEnd<br/>"]
    q277{"CompositeAlternatives__Basic_3 (277)<br/>Basic<br/><br/>dec=32"}

    q46 --> q271
    q271 -.->|"[CompositeGroup]"| q277
    q272 -->|"tok(&quot;|&quot;)"| q273
    q273 -.->|"[CompositeGroup]"| q274
    q274 --> q275
    q275 --> q272
    q275 --> q276
    q276 --> q47
    q277 --> q272
    q277 --> q276
```

## CompositeGroup

```mermaid
flowchart TD
    q48(["CompositeGroup__Start (48)<br/>RuleStart"])
    q49(["CompositeGroup__Stop (49)<br/>RuleStop"])
    q278["CompositeGroup__Basic_0 (278)<br/>Basic<br/>"]
    q279["CompositeGroup__Basic_1 (279)<br/>Basic<br/>"]
    q280["CompositeGroup__Basic_2 (280)<br/>Basic<br/>"]
    q281{"CompositeGroup__LoopBack (281)<br/>LoopBack<br/><br/>dec=33"}
    q282["CompositeGroup__LoopEnd (282)<br/>LoopEnd<br/>"]
    q283{"CompositeGroup__Basic_3 (283)<br/>Basic<br/><br/>dec=34"}

    q48 --> q278
    q278 -.->|"[CompositeElement]"| q283
    q279 -.->|"[CompositeElement]"| q280
    q280 --> q281
    q281 --> q279
    q281 --> q282
    q282 --> q49
    q283 --> q279
    q283 --> q282
```

## CompositeElement

```mermaid
flowchart TD
    q50(["CompositeElement__Start (50)<br/>RuleStart"])
    q51(["CompositeElement__Stop (51)<br/>RuleStop"])
    q284["CompositeElement__Basic_0 (284)<br/>Basic<br/>"]
    q285["CompositeElement__Basic_1 (285)<br/>Basic<br/>"]
    q286["CompositeElement__Basic_2 (286)<br/>Basic<br/>"]
    q287["CompositeElement__Basic_3 (287)<br/>Basic<br/>"]
    q288["CompositeElement_LeftParen (288)<br/>Basic<br/>"]
    q289["CompositeElement__Basic_4 (289)<br/>Basic<br/>"]
    q290["CompositeElement_RightParen (290)<br/>Basic<br/>"]
    q291["CompositeElement__Basic_5 (291)<br/>Basic<br/>"]
    q292{"CompositeElement__Basic_6 (292)<br/>Basic<br/><br/>dec=35"}
    q293["CompositeElement__BlockEnd_0 (293)<br/>BlockEnd<br/>"]
    q294["CompositeElement_Cardinality_Asterisk (294)<br/>Basic<br/>"]
    q295["CompositeElement__Basic_7 (295)<br/>Basic<br/>"]
    q296["CompositeElement_Cardinality_Plus (296)<br/>Basic<br/>"]
    q297["CompositeElement__Basic_8 (297)<br/>Basic<br/>"]
    q298["CompositeElement_Cardinality_Question (298)<br/>Basic<br/>"]
    q299["CompositeElement__Basic_9 (299)<br/>Basic<br/>"]
    q300{"CompositeElement__Basic_10 (300)<br/>Basic<br/><br/>dec=36"}
    q301["CompositeElement__BlockEnd_1 (301)<br/>BlockEnd<br/>"]
    q302{"CompositeElement__Basic_11 (302)<br/>Basic<br/><br/>dec=37"}

    q50 --> q292
    q284 -.->|"[Keyword]"| q285
    q285 --> q293
    q286 -.->|"[RuleCall]"| q287
    q287 --> q293
    q288 -->|"tok(&quot;(&quot;)"| q289
    q289 -.->|"[CompositeAlternatives]"| q290
    q290 -->|"tok(&quot;)&quot;)"| q291
    q291 --> q293
    q292 --> q284
    q292 --> q286
    q292 --> q288
    q293 --> q302
    q294 -->|"tok(&quot;*&quot;)"| q295
    q295 --> q301
    q296 -->|"tok(&quot;+&quot;)"| q297
    q297 --> q301
    q298 -->|"tok(&quot;?&quot;)"| q299
    q299 --> q301
    q300 --> q294
    q300 --> q296
    q300 --> q298
    q301 --> q51
    q302 --> q300
    q302 --> q301
```

