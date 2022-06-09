package mngAdmin

var SelectBoardCnt = `
SELECT Count(*) as TOTAL_COUNT
FROM sys_boards 
WHERE 
	TITLE LIKE '%#{search}%'
`

var SelectBoardList = `
SELECT 
		BOARD_ID
		,B_KIND
		,BOARD_TYPE
		,TITLE
		,USE_YN
		,START_DATE
		,END_DATE
		,DATE_FORMAT(REG_DATE,'%Y-%m-%d') as regDate
		,DATE_FORMAT(MOD_DATE,'%Y-%m-%d') as modDate
FROM sys_boards 
WHERE 
	TITLE LIKE '%#{search}%'
`

var InsertBoard = `
INSERT INTO sys_boards (
	B_KIND
	,USE_YN
	,MAIN_YN
	,BOARD_TYPE
	,START_DATE
	,END_DATE
	,TITLE 
	,LINK_URL 
	,CONTENT
	,REG_DATE
) VALUES (
	'#{bKind}'
	,'#{useYn}'
	,'#{mainYn}'
	,'#{boardType}'
	,'#{startDate}'
	,'#{endDate}'
	,'#{title}'
	,IF('#{link}'='',NULL,'#{link}')
	,IF('#{content}'='',NULL,'#{content}')
	,DATE_FORMAT(now(),'%Y%m%d%H%i%s')
)
`

var SelectBoardInfo = `
SELECT 
		BOARD_ID
		,B_KIND
		,BOARD_TYPE
		,TITLE
		,IFNULL(LINK_URL,'') as linkUrl
		,IFNULL(CONTENT,'') as Content
		,USE_YN
		,MAIN_YN
		,START_DATE
		,END_DATE
		,REG_DATE
		,MOD_DATE
FROM sys_boards 
WHERE 
	BOARD_ID = '#{boardId}'
`

var UpdateBoardInfo = `
UPDATE sys_boards 
SET 
	B_KIND = '#{bKind}'
	,BOARD_TYPE = '#{boardType}'
	,TITLE = '#{title}'
	,LINK_URL = '#{link}'
	,CONTENT = '#{content}'
	,USE_YN = '#{useYn}'
	,MAIN_YN = '#{mainYn}'
	,START_DATE = '#{startDate}'
	,END_DATE = '#{endDate}'
	,MOD_DATE = DATE_FORMAT(now(), '%Y%m%d%H%i%s')
WHERE 
	BOARD_ID = '#{boardId}'
`
var SelectContentList = `
SELECT 
	CONTENT_ID
	,TYPE
	,SITE_NAME
	,TITLE
	,USE_YN
	,START_DATE
	,END_DATE
FROM b_contents
WHERE
	TITLE LIKE '%#{search}%'
`

var SelectContentCnt = `
SELECT 
	Count(*) as TOTAL_COUNT
FROM b_contents
WHERE
	TITLE LIKE '%#{search}%'
`

var SelectContentInfo = `
SELECT 
	CONTENT_ID
	,TYPE
	,SITE_NAME
	,TITLE
	,USE_YN
	,DESCRIPTION
	,URL
	,IMAGE_URL
	,VIDEO_URL
	,START_DATE
	,END_DATE
FROM b_contents
WHERE
	CONTENT_ID = '#{contentId}'
`

var UpdateContentInfo = `
UPDATE b_contents 
SET
	TYPE = '#{type}'
	,SITE_NAME = '#{siteName}'
	,USE_YN = '#{useYn}'	
	,START_DATE = '#{startDate}'
	,END_DATE = '#{endDate}'
	,TITLE = '#{title}'
	,URL = '#{url}'
	,IMAGE_URL = '#{imageUrl}'
	,VIDEO_URL = '#{videoUrl}'
	,DESCRIPTION = '#{description}'
	,MOD_DATE = DATE_FORMAT(now(),'%Y-%m-%d')
WHERE
	CONTENT_ID = '#{contentId}'
`
var SelectBannerList = `
SELECT 
	BANNER_ID
	,B_KIND
	,BANNER_TYPE
	,TITLE
	,USE_YN
	,START_DATE
	,END_DATE
	,DATE_FORMAT(IFNULL(MOD_DATE,REG_DATE),'%Y-%m-%d') as modDate
FROM b_banners 
WHERE 
	TITLE LIKE '%#{search}%'
`

var SelectBannerCnt = `
SELECT Count(*) as TOTAL_COUNT
FROM b_banners 
WHERE 
	TITLE LIKE '%#{search}%'
`

var InsertBanner = `
INSERT INTO b_banners (
	TITLE 
	,URL 
	,IMAGE
	,B_KIND
	,USE_YN
	,BANNER_TYPE
	,START_DATE
	,END_DATE
	,START_TIME
	,END_TIME
	,REG_DATE
) VALUES (
	'#{title}'
	,IF('#{link}'='',NULL,'#{link}')
	,IF('#{image}'='',NULL,'#{image}')
	,'#{bKind}'
	,'#{useYn}'
	,'#{bannerType}'
	,'#{startDate}'
	,'#{endDate}'
	,'#{startTime}'
	,'#{endTime}'
	,DATE_FORMAT(now(),'%Y%m%d%H%i%s')
)
`

var SelectBannerInfo = `
SELECT 
	BANNER_ID
	,B_KIND
	,BANNER_TYPE
	,TITLE
	,USE_YN
	,IFNULL(URL,'') as url
	,IFNULL(IMAGE,'') as image
	,START_DATE
	,END_DATE
	,START_TIME
	,END_TIME
	,DATE_FORMAT(IFNULL(MOD_DATE,REG_DATE),'%Y-%m-%d') as modDate
FROM b_banners 
WHERE 
	BANNER_ID = '#{bannerId}'
`

var UpdateBannerInfo = `
UPDATE b_banners
SET 
	B_KIND = '#{bKind}'
	,BANNER_TYPE = '#{bannerType}'
	,USE_YN = '#{useYn}'
	,START_DATE = '#{startDate}'
	,END_DATE = '#{endDate}'
	,START_TIME = '#{startTime}'
	,END_TIME = '#{endTime}'
	,TITLE = '#{title}'
	,URL = '#{link}'
	,IMAGE = '#{image}'
	,MOD_DATE = DATE_FORMAT(now(), '%Y%m%d%H%i%s')
WHERE 
	BANNER_ID = '#{bannerId}'
`
