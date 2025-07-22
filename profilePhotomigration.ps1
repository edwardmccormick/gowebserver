# Example: List of people with ID and Profile S3Key
$people = @(
    @{ ID = 1; S3Key = "example/man_square2.jpg" },
    @{ ID = 2; S3Key = "example/pexels-kelvin809-810775.jpg" },
    @{ ID = 3; S3Key = "example/turdfurguson.jpg" },
    @{ ID = 4; S3Key = "example/President-Donald-Trump-Official-Presidential-Portrait.png" },
    @{ ID = 5; S3Key = "example/tedProfilePicture.jpg" },
    @{ ID = 6; S3Key = "example/image-Scarlett-Johansson-Images.jpg" },
    @{ ID = 7; S3Key = "example/Lenna.jpg" },
    @{ ID = 8; S3Key = "example/natalie.jpg" }
)

$bucket = "urmid-images"

foreach ($person in $people) {
    $source = "$bucket/$($person.S3Key)"
    $destination = "$bucket/$($person.ID)/profile"

    # Copy the object to the new key
    aws s3 cp "s3://$source" "s3://$destination"

    # Optionally, delete the original object
    # aws s3 rm "s3://$source"
}