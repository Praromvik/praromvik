## API
```
type Course struct {
    CourseId           string    `json:"_id" bson:"_id"`
    Title              string    `json:"title" bson:"title"`
    Description        string    `json:"description" bson:"description"`
    Instructors        []string  `json:"instructors" bson:"instructors"`
    Moderators         []string  `json:"moderators" bson:"moderators"`
    StartDate          date.Date `json:"startDate" bson:"startDate"`
    EndDate            date.Date `json:"endDate" bson:"endDate"`
    Duration           int       `json:"duration" bson:"duration"` // Duration in week
    Capacity           int       `json:"capacity" bson:"capacity"`
    Students           []string  `json:"students" bson:"students"`
    Price              int       `json:"price" bson:"price"`
    Image              []byte    `json:"image" bson:"-"`
}
```

```
type Lesson struct {
    LessonID     string       `json:"_id" bson:"_id"`
    Title        string       `json:"title" bson:"title"`
    Contents     []ContentRef `json:"contents" bson:"contents"`
}

type ContentRef struct {
    ID      string `json:"id" bson:"id"`
    Title   string `json:"title" bson:"title"`
}

type Content struct {
    ContentID    string `json:"_id" bson:"_id"`
    Title        string `json:"title" bson:"title"`
    Type         string `json:"type" bson:"type"` // video, resource, quiz, lab
    Data         []byte `json:"data" bson:"data"`
    LessonRef    UUID   `json:"lessonRef" bson:"lessonRef"`
}
```


## A concrete Example
```
Course name: Mastering golang
<Lessons>:
Introduction:
    Why Go?
    Installation
    <resources>
    <quiz>
Concurrency in Go:
    channels
    mutex
    <lab1>
    <resources>
    go routines
    <lab2>
Conclusion
    thank you
```

For the `lab1` under 'Concurrency in Go', url will be like:
http://praromvik.com/course/mastering_golang/concurrency_lab_1/
Intentionally cutting off the lesson id part to keep the url shorter.

---

We have divided a single course into 2 different databases. One for holding the metadata, & one for holding the actual information.
Again the actual information is divided into two collections. One will be called once, & the other one will be called several times (with specific contentID) from the frontend.
### `praromvik.courses` collection
This will hold the metadata information only.
```
{
    "_id": "mastering_golang",
    "title": "Mastering Golang",
    "instructors": []string{<some_user_ids>},
    "students": []string{<some_user_ids>}
}
```
Removed the lessons slice above to avoid multiple reference update on create/delete calls.

### `mastering_golang.lessons` collection
The _id format is `<4_digit_serial_number>_<epoch_time>`
```
{
    "_id": "<0001_1715524983>",
    "title": "Introduction",
    "contents": []string{<_ids>}
},
{
    "_id": "<0002_1715527643>",
    "title": "Concurrency in Go",
    "contents": []string{"channel", "concurrency_lab_1", "resouce_from_doc"}
}
```

### `mastering_golang.contents` collection
```
{
    "_id": "channel",
    "title": "Channel",
    "type": "Video",
    "lessonRef": "0002_1715527643"
},
{
    "_id": "concurrency_lab_1",
    "title": "Hands-on the Go channel",
    "type": "Lab",
    "lessonRef": "0002_1715527643"
},
{
    "_id": "resouce_from_doc",
    "title": "Additional learing resouces",
    "type": "Written",
    "lessonRef": "0002_1715527643"
},
```

So the `RESTRICTIONS` are:

i) course ids have to be unique.

ii) content ids within a course have to be unique.

iii) both course & content ids have to be <= 24 characters, as they will be used as _id in mongodb.
