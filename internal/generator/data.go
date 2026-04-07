package generator

import (
	"fmt"
	"math/rand"
)

// ---- GENDERED NAMES ----

var maleNames = []string{
	"James", "Robert", "John", "Michael", "David", "William", "Richard", "Joseph",
	"Thomas", "Christopher", "Charles", "Daniel", "Matthew", "Anthony", "Mark",
	"Donald", "Steven", "Paul", "Andrew", "Joshua", "Kenneth", "Kevin", "Brian",
	"George", "Timothy", "Ronald", "Edward", "Jason", "Jeffrey", "Ryan",
	"Jacob", "Gary", "Nicholas", "Eric", "Jonathan", "Stephen", "Larry", "Justin",
	"Scott", "Brandon", "Benjamin", "Samuel", "Raymond", "Gregory", "Frank",
	"Alexander", "Patrick", "Jack", "Dennis", "Jerry", "Tyler", "Aaron", "Jose",
	"Nathan", "Henry", "Peter", "Douglas", "Zachary", "Kyle", "Noah", "Ethan",
	"Jeremy", "Walter", "Christian", "Keith", "Roger", "Terry", "Austin", "Sean",
	"Gerald", "Carl", "Harold", "Dylan", "Arthur", "Lawrence", "Jordan", "Jesse",
	"Bryan", "Billy", "Bruce", "Gabriel", "Joe", "Logan", "Albert", "Willie",
	"Alan", "Vincent", "Philip", "Bobby", "Johnny", "Bradley", "Roy", "Ralph",
	"Randy", "Wayne", "Howard", "Adam", "Carlos", "Marcus", "Oscar", "Liam",
	"Mason", "Lucas", "Oliver", "Elijah", "Aiden", "Sebastian", "Caleb", "Owen",
	"Connor", "Jayden", "Isaiah", "Luke", "Adrian", "Eli", "Max", "Nolan",
	"Parker", "Leo", "Miles", "Dominic", "Grant", "Colton", "Levi", "Chase",
	"Blake", "Declan", "Gavin", "Wesley", "Cole", "Tristan", "Micah", "Damian",
	"Roman", "Ivan", "Axel", "Harrison", "Jonah", "Rafael", "Spencer", "Elliott",
	"Victor", "Shane", "Graham", "Trevor", "Jared", "Preston", "Derek", "Bryce",
	"Cody", "Travis", "Mitchell", "Dustin", "Craig", "Kirk", "Ross", "Brent",
	"Todd", "Neil", "Troy", "Corey", "Chad", "Kurt", "Lance", "Darren", "Wade",
	"Martin", "Russell", "Eugene", "Ernest", "Phillip", "Francis", "Clarence",
	"Louis", "Stanley", "Leonard", "Dale", "Manuel", "Rodney", "Curtis",
}

var femaleNames = []string{
	"Mary", "Patricia", "Jennifer", "Linda", "Elizabeth", "Barbara", "Susan", "Jessica",
	"Sarah", "Karen", "Lisa", "Nancy", "Betty", "Margaret", "Sandra", "Ashley",
	"Kimberly", "Emily", "Donna", "Michelle", "Carol", "Amanda", "Dorothy", "Melissa",
	"Deborah", "Stephanie", "Rebecca", "Sharon", "Laura", "Cynthia", "Kathleen", "Amy",
	"Angela", "Shirley", "Anna", "Brenda", "Pamela", "Emma", "Nicole", "Helen",
	"Samantha", "Katherine", "Christine", "Debra", "Rachel", "Carolyn", "Janet", "Catherine",
	"Maria", "Heather", "Diane", "Ruth", "Julie", "Olivia", "Joyce", "Virginia",
	"Victoria", "Kelly", "Lauren", "Christina", "Joan", "Evelyn", "Judith", "Megan",
	"Andrea", "Cheryl", "Hannah", "Jacqueline", "Martha", "Gloria", "Teresa", "Ann",
	"Sara", "Madison", "Frances", "Kathryn", "Janice", "Jean", "Abigail", "Alice",
	"Julia", "Judy", "Sophia", "Grace", "Denise", "Amber", "Doris", "Marilyn",
	"Danielle", "Beverly", "Isabella", "Theresa", "Diana", "Natalie", "Brittany", "Charlotte",
	"Marie", "Kayla", "Alexis", "Lori", "Ella", "Mia", "Ava", "Harper",
	"Aria", "Scarlett", "Penelope", "Layla", "Chloe", "Riley", "Zoey", "Nora",
	"Lily", "Eleanor", "Hazel", "Violet", "Aurora", "Savannah", "Audrey", "Brooklyn",
	"Bella", "Claire", "Skylar", "Lucy", "Paisley", "Everly", "Stella", "Caroline",
	"Maya", "Naomi", "Elena", "Piper", "Lydia", "Alexa", "Josephine", "Allison",
	"Madeline", "Peyton", "Kennedy", "Ruby", "Ivy", "Gabriella", "Gianna", "Mackenzie",
	"Autumn", "Quinn", "Bailey", "Serenity", "Aubrey", "Morgan", "Taylor", "Leah",
	"Faith", "Melanie", "Paige", "Brooke", "Adriana", "Jenna", "Marissa", "Tiffany",
	"Courtney", "Vanessa", "Lindsey", "Erica", "Kristin", "Dana", "Crystal", "Tara",
	"Alicia", "Holly", "Meredith", "Cassandra", "Wendy", "Kelsey", "Tamara", "Monica",
}

var allLastNames = []string{
	"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis",
	"Rodriguez", "Martinez", "Hernandez", "Lopez", "Gonzalez", "Wilson", "Anderson",
	"Thomas", "Taylor", "Moore", "Jackson", "Martin", "Lee", "Perez", "Thompson",
	"White", "Harris", "Sanchez", "Clark", "Ramirez", "Lewis", "Robinson", "Walker",
	"Young", "Allen", "King", "Wright", "Scott", "Torres", "Nguyen", "Hill",
	"Flores", "Green", "Adams", "Nelson", "Baker", "Hall", "Rivera", "Campbell",
	"Mitchell", "Carter", "Roberts", "Gomez", "Phillips", "Evans", "Turner", "Diaz",
	"Parker", "Cruz", "Edwards", "Collins", "Reyes", "Stewart", "Morris", "Morales",
	"Murphy", "Cook", "Rogers", "Gutierrez", "Ortiz", "Morgan", "Cooper", "Peterson",
	"Bailey", "Reed", "Kelly", "Howard", "Ramos", "Kim", "Cox", "Ward",
	"Richardson", "Watson", "Brooks", "Chavez", "Wood", "James", "Bennett", "Gray",
	"Mendoza", "Ruiz", "Hughes", "Price", "Alvarez", "Castillo", "Sanders", "Patel",
	"Myers", "Long", "Ross", "Foster", "Jimenez", "Powell", "Jenkins", "Perry",
	"Russell", "Sullivan", "Bell", "Coleman", "Butler", "Henderson", "Barnes", "Fisher",
	"Vasquez", "Simmons", "Graham", "Murray", "Ford", "Stone", "Hunt", "Dunn",
	"Dean", "Gibson", "Wells", "Palmer", "Burton", "Fox", "Webb", "Porter",
	"Burke", "Walsh", "Hammond", "Tucker", "Garrett", "Pearson", "Ferguson", "Hawkins",
	"Spencer", "Boyd", "Mason", "Arnold", "Wagner", "Lynch", "Drake", "Fleming",
	"Harrison", "Mccarthy", "Keller", "Abbott", "Chambers", "Caldwell", "Maxwell",
	"Dixon", "Carr", "Watts", "Hart", "Day", "Lawson", "Knight", "Barker",
	"Shelton", "Weaver", "Barton", "Burns", "Payne", "Mcdonald", "Owens", "Bates",
	"Freeman", "Holt", "French", "Marsh", "Stephens", "Crawford", "Riley", "Duffy",
	"Hoffman", "Olson", "Hansen", "Fernandez", "Garza", "Harvey", "Burton", "Wolfe",
	"Bishop", "Mccoy", "Howell", "Larson", "Brewer", "Pratt", "Hines", "Gallagher",
	"Roberson", "Potts", "Dillon", "Saunders", "Frost", "Donovan", "Sutton", "Malone",
}

// ---- LOCATION DATA (US only, coherent state→city→zip→area code) ----

type stateData struct {
	Code      string
	AreaCodes []string
	Cities    []cityData
}

type cityData struct {
	Name      string
	ZipPrefix int // first 3 digits of zip code
}

var usStates = []stateData{
	{"CA", []string{"213", "310", "415", "408", "510", "619", "714", "818", "858", "916"}, []cityData{
		{"Los Angeles", 900}, {"San Francisco", 941}, {"San Diego", 921}, {"San Jose", 951},
		{"Sacramento", 958}, {"Oakland", 946}, {"Fresno", 937}, {"Long Beach", 908},
		{"Santa Barbara", 931}, {"Pasadena", 911}, {"Berkeley", 947}, {"Irvine", 926},
	}},
	{"TX", []string{"210", "214", "281", "512", "713", "817", "832", "903", "915", "972"}, []cityData{
		{"Houston", 770}, {"Dallas", 752}, {"Austin", 787}, {"San Antonio", 782},
		{"Fort Worth", 761}, {"El Paso", 799}, {"Arlington", 760}, {"Plano", 750},
		{"Corpus Christi", 784}, {"Lubbock", 794}, {"Laredo", 780}, {"Amarillo", 791},
	}},
	{"NY", []string{"212", "315", "347", "516", "518", "585", "607", "631", "646", "716", "718", "914"}, []cityData{
		{"New York", 100}, {"Brooklyn", 112}, {"Buffalo", 142}, {"Rochester", 146},
		{"Albany", 122}, {"Syracuse", 132}, {"Yonkers", 107}, {"White Plains", 106},
		{"Ithaca", 148}, {"Schenectady", 123},
	}},
	{"FL", []string{"239", "305", "321", "352", "386", "407", "561", "727", "786", "813", "850", "904", "941", "954"}, []cityData{
		{"Miami", 331}, {"Orlando", 328}, {"Tampa", 336}, {"Jacksonville", 322},
		{"St. Petersburg", 337}, {"Fort Lauderdale", 333}, {"Tallahassee", 323},
		{"Gainesville", 326}, {"Sarasota", 342}, {"Naples", 341},
	}},
	{"IL", []string{"217", "309", "312", "618", "630", "708", "773", "815", "847"}, []cityData{
		{"Chicago", 606}, {"Aurora", 605}, {"Springfield", 627}, {"Peoria", 616},
		{"Rockford", 611}, {"Naperville", 605}, {"Joliet", 604}, {"Elgin", 601},
	}},
	{"PA", []string{"215", "267", "412", "484", "570", "610", "717", "814"}, []cityData{
		{"Philadelphia", 191}, {"Pittsburgh", 152}, {"Allentown", 181}, {"Erie", 165},
		{"Reading", 196}, {"Lancaster", 176}, {"Harrisburg", 171}, {"Bethlehem", 180},
	}},
	{"OH", []string{"216", "234", "330", "419", "440", "513", "614", "740", "937"}, []cityData{
		{"Columbus", 432}, {"Cleveland", 441}, {"Cincinnati", 452}, {"Toledo", 436},
		{"Akron", 443}, {"Dayton", 454}, {"Canton", 447}, {"Youngstown", 445},
	}},
	{"GA", []string{"229", "404", "478", "678", "706", "770", "912"}, []cityData{
		{"Atlanta", 303}, {"Augusta", 309}, {"Savannah", 314}, {"Columbus", 319},
		{"Athens", 306}, {"Macon", 312}, {"Roswell", 300}, {"Albany", 317},
	}},
	{"NC", []string{"252", "336", "704", "828", "910", "919", "980"}, []cityData{
		{"Charlotte", 282}, {"Raleigh", 276}, {"Durham", 277}, {"Greensboro", 274},
		{"Winston-Salem", 271}, {"Fayetteville", 283}, {"Cary", 275}, {"Wilmington", 284},
	}},
	{"MI", []string{"231", "248", "269", "313", "517", "586", "616", "734", "810", "947"}, []cityData{
		{"Detroit", 482}, {"Grand Rapids", 495}, {"Ann Arbor", 481}, {"Lansing", 489},
		{"Flint", 485}, {"Kalamazoo", 490}, {"Sterling Heights", 483}, {"Warren", 480},
	}},
	{"WA", []string{"206", "253", "360", "425", "509"}, []cityData{
		{"Seattle", 981}, {"Tacoma", 984}, {"Spokane", 992}, {"Vancouver", 986},
		{"Bellevue", 980}, {"Olympia", 985}, {"Everett", 982}, {"Redmond", 980},
	}},
	{"MA", []string{"339", "351", "413", "508", "617", "774", "781", "857", "978"}, []cityData{
		{"Boston", 21}, {"Worcester", 16}, {"Springfield", 11}, {"Cambridge", 21},
		{"Lowell", 18}, {"Salem", 19}, {"Quincy", 21}, {"Newton", 24},
	}},
	{"CO", []string{"303", "719", "720", "970"}, []cityData{
		{"Denver", 802}, {"Colorado Springs", 809}, {"Aurora", 800}, {"Fort Collins", 805},
		{"Boulder", 803}, {"Pueblo", 810}, {"Lakewood", 802}, {"Arvada", 800},
	}},
	{"AZ", []string{"480", "520", "602", "623", "928"}, []cityData{
		{"Phoenix", 850}, {"Tucson", 857}, {"Mesa", 852}, {"Scottsdale", 852},
		{"Chandler", 852}, {"Tempe", 852}, {"Glendale", 853}, {"Flagstaff", 860},
	}},
	{"VA", []string{"276", "434", "540", "571", "703", "757", "804"}, []cityData{
		{"Richmond", 232}, {"Virginia Beach", 234}, {"Norfolk", 235}, {"Arlington", 222},
		{"Alexandria", 223}, {"Chesapeake", 233}, {"Roanoke", 240}, {"Newport News", 236},
	}},
	{"TN", []string{"423", "615", "629", "731", "865", "901", "931"}, []cityData{
		{"Nashville", 372}, {"Memphis", 381}, {"Knoxville", 379}, {"Chattanooga", 374},
		{"Clarksville", 370}, {"Murfreesboro", 371}, {"Franklin", 370}, {"Jackson", 383},
	}},
	{"OR", []string{"458", "503", "541", "971"}, []cityData{
		{"Portland", 972}, {"Salem", 973}, {"Eugene", 974}, {"Bend", 977},
		{"Medford", 975}, {"Corvallis", 973}, {"Hillsboro", 971}, {"Gresham", 970},
	}},
	{"MN", []string{"218", "320", "507", "612", "651", "763", "952"}, []cityData{
		{"Minneapolis", 554}, {"St. Paul", 551}, {"Rochester", 559}, {"Duluth", 558},
		{"Bloomington", 554}, {"Plymouth", 553}, {"Woodbury", 551}, {"Eagan", 551},
	}},
	{"NV", []string{"702", "725", "775"}, []cityData{
		{"Las Vegas", 891}, {"Reno", 895}, {"Henderson", 890}, {"North Las Vegas", 890},
		{"Sparks", 894}, {"Carson City", 897},
	}},
	{"NJ", []string{"201", "551", "609", "732", "848", "856", "862", "908", "973"}, []cityData{
		{"Newark", 71}, {"Jersey City", 73}, {"Trenton", 86}, {"Princeton", 85},
		{"Camden", 81}, {"Hoboken", 70}, {"Morristown", 79}, {"Elizabeth", 72},
	}},
}

// ---- STREETS ----

var streetNames = []string{
	"Main", "Oak", "Pine", "Maple", "Cedar", "Elm", "Washington", "Lake",
	"Hill", "Park", "Walnut", "Sunset", "Jackson", "Church", "River",
	"Lincoln", "Spring", "Highland", "Franklin", "Jefferson", "Adams", "Madison",
	"Willow", "Meadow", "Forest", "Bridge", "Cherry", "Vine", "Pleasant",
	"Prospect", "Center", "School", "Union", "Mill", "Market", "Broad",
	"Chestnut", "Liberty", "Pearl", "Academy", "Grove", "College", "Poplar",
	"Hickory", "Laurel", "Birch", "Spruce", "Dogwood", "Magnolia", "Sycamore",
	"Aspen", "Beech", "Mulberry", "Orchard", "Ridge", "Valley", "Summit",
	"Lakeview", "Fairway", "Woodland", "Creekside", "Riverside", "Bayview",
}

var streetSuffixes = []string{
	"St", "Ave", "Blvd", "Dr", "Ln", "Way", "Ct", "Rd", "Pl", "Ter", "Cir",
}

// ---- SAFE DOMAINS (RFC 2606 reserved) ----

var safePersonalDomains = []string{
	"gmail.com", "yahoo.com", "hotmail.com", "outlook.com", "protonmail.com",
	"icloud.com", "fastmail.com", "mail.com",
}

// ---- HELPER FUNCTIONS ----

// pickGendered returns a coherent first name, last name, and prefix.
func pickGendered(rng *rand.Rand) (first, last, prefix string, isMale bool) {
	last = allLastNames[rng.Intn(len(allLastNames))]

	if rng.Intn(2) == 0 {
		isMale = true
		first = maleNames[rng.Intn(len(maleNames))]
		if rng.Intn(10) == 0 {
			prefix = "Dr."
		} else {
			prefix = "Mr."
		}
	} else {
		isMale = false
		first = femaleNames[rng.Intn(len(femaleNames))]
		switch {
		case rng.Intn(10) == 0:
			prefix = "Dr."
		case rng.Intn(2) == 0:
			prefix = "Mrs."
		default:
			prefix = "Ms."
		}
	}
	return
}

// pickLocation returns a coherent US state code, city, zip, and area code.
func pickLocation(rng *rand.Rand) (stateCode, city, zip, areaCode string) {
	st := usStates[rng.Intn(len(usStates))]
	ct := st.Cities[rng.Intn(len(st.Cities))]
	ac := st.AreaCodes[rng.Intn(len(st.AreaCodes))]
	zip = fmt.Sprintf("%03d%02d", ct.ZipPrefix, rng.Intn(100))
	return st.Code, ct.Name, zip, ac
}

// pickStreet returns a street address like "1247 Oak Ave".
func pickStreet(rng *rand.Rand) string {
	num := 100 + rng.Intn(9900)
	street := streetNames[rng.Intn(len(streetNames))]
	suffix := streetSuffixes[rng.Intn(len(streetSuffixes))]
	return fmt.Sprintf("%d %s %s", num, street, suffix)
}

// genUUID generates a v4-format UUID from the given rng.
func genUUID(rng *rand.Rand) string {
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		rng.Uint32(),
		rng.Intn(0xFFFF),
		0x4000|rng.Intn(0x0FFF),
		0x8000|rng.Intn(0x3FFF),
		rng.Int63n(0xFFFFFFFFFFFF),
	)
}
