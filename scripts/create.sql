IF OBJECT_ID(N'[dbo].[Department]', N'U') IS null
    BEGIN
        CREATE TABLE [dbo].[Department]
        (
            [ID]   [int] IDENTITY (1,1) NOT NULL PRIMARY KEY,
            [NAME] [nvarchar](MAX)      NOT NULL,
        )
        INSERT INTO [dbo].[Department]
        VALUES ('HR'),
               ('IT'),
               ('DESIGN')
    END

IF OBJECT_ID(N'[dbo].[Employee]', N'U') IS null
    BEGIN
        CREATE TABLE [dbo].[Employee]
        (
            [ID]            [int] IDENTITY (1,1) NOT NULL PRIMARY KEY,
            [FIRST_NAME]    [nvarchar](MAX)      NOT NULL,
            [LAST_NAME]     [nvarchar](MAX)      NOT NULL,
            [USERNAME]      [nvarchar](MAX)      NOT NULL,
            [PASSWORD]      [varbinary](160)     NOT NULL,
            [EMAIL]         [nvarchar](MAX)      NOT NULL,
            [DOB]           [date]               NOT NULL,
            [POSITION]      [nvarchar](MAX)      NOT NULL,
            [DEPARTMENT_ID] [int] FOREIGN KEY REFERENCES [dbo].[Department] (ID),
        )
        INSERT INTO [dbo].[Employee]
        VALUES ('Connor', 'MacLeod', 'highlander', PWDENCRYPT('whoWantsToLiveForEver'),'high_lander@gmail.com',
                '29-aug-1986', 'Head of Tech', 2)
    END