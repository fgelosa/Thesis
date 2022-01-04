import com.owlike.genson.annotation.JsonProperty;
import org.hyperledger.fabric.contract.annotation.DataType;
import org.hyperledger.fabric.contract.annotation.Property;

@DataType
public class Hash {

    @Property()
    private final String hashString;

    @Property()
    private final String date;

    @Property()
    private final String user;

    public String getHashString(){
        return hashString;
    }

    public String getDate() {
        return date;
    }

    public String getUser() {
        return user;
    }

    public Hash(@JsonProperty("hash") final String hashString,@JsonProperty("date") final String date,@JsonProperty("user") final String user){
        this.hashString = hashString;
        this.date = date;
        this.user = user;
    }
}
